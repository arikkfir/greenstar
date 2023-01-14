package resolver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arikkfir/greenstar/api/model"
	"github.com/rs/zerolog/log"
	"github.com/rueian/rueidis"
	"github.com/xuri/excelize/v2"
	"regexp"
	"time"
)

func (r *transactionsXLSXProcessor) Run(ctx context.Context) {
	cmd := r.Redis.B().Subscribe().Channel(processXLSXChannelName).Build()
	err := r.Redis.Receive(ctx, cmd, func(msg rueidis.PubSubMessage) {
		var data []byte
		if err := json.Unmarshal([]byte(msg.Message), &data); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to process XLSX message")
		} else if err := r.processFile(ctx, data); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to process message")
		}
	})
	if err != nil && !errors.Is(err, rueidis.ErrClosing) {
		log.Ctx(ctx).Fatal().Err(err).Msg("Failed to subscribe to Redis queue")
	}
}

func (r *transactionsXLSXProcessor) processFile(ctx context.Context, data []byte) error {
	const (
		transactionsSheetName        = "תנועות עו\"ש"
		holdingTransactionsSheetName = "תנועות זמניות בהמתנה"
	)
	var (
		expectedCellContents = map[string]string{
			"F2": "תנועות בחשבון עו\"ש",
			"F3": "תנועות עו\"ש",
			"I4": "חשבון",
			"I5": "מועד הפקת הדוח",
			"A6": "יתרה משוערכת(₪)",
			"B6": "זכות(₪)",
			"C6": "חובה(₪)",
			"D6": "תיאור פעולה",
			"E6": "אסמכתא",
			"H6": "תאריך ערך",
			"J6": "תאריך",
		}
		accountFormatRE = regexp.MustCompile("^([0-9]{2})-([0-9]{3})-([0-9]{6})$")
	)

	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to open uploaded Excel file: %w", err)
	} else if f.SheetCount != 2 {
		return fmt.Errorf("expected 2 sheets, got %d", f.SheetCount)
	}

	transactionsSheet := f.WorkBook.Sheets.Sheet[0]
	if transactionsSheet.Name != transactionsSheetName {
		return fmt.Errorf("expected sheet name to be '%s', got '%s'", transactionsSheetName, transactionsSheet.Name)
	}

	holdingTransactionsSheet := f.WorkBook.Sheets.Sheet[1]
	if holdingTransactionsSheet.Name != holdingTransactionsSheetName {
		return fmt.Errorf("expected sheet name to be '%s', got '%s'", holdingTransactionsSheetName, holdingTransactionsSheet.Name)
	}

	for cell, expectedContent := range expectedCellContents {
		actualContent, err := f.GetCellValue(transactionsSheet.Name, cell)
		if err != nil {
			return fmt.Errorf("failed to get cell value for cell '%s': %w", cell, err)
		} else if actualContent != expectedContent {
			return fmt.Errorf("expected cell '%s' to have content '%s', got '%s'", cell, expectedContent, actualContent)
		}
	}

	var branch, account string
	a4Value, err := f.GetCellValue(transactionsSheet.Name, "A4", excelize.Options{RawCellValue: true})
	if err != nil {
		return fmt.Errorf("failed to get cell value for cell 'A4': %w", err)
	} else if match := accountFormatRE.FindStringSubmatch(a4Value); match == nil {
		return fmt.Errorf("expected cell 'A4' to have content matching '%s', got '%s'", accountFormatRE.String(), a4Value)
	} else if bank := match[1]; bank != "04" {
		// Bank Yahav is 04
		return fmt.Errorf("expected cell 'A4' to have bank '04', got '%s'", bank)
	} else {
		branch = match[2]
		account = match[3]
	}

	rows, err := f.GetRows(transactionsSheet.Name, excelize.Options{})
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet '%s': %w", transactionsSheet.Name, err)
	} else if len(rows) < 7 {
		// no transactions
		return nil
	}

	var inputs []*model.NewTransaction
	for _, row := range rows[6:] {
		if len(row) == 10 {
			credit, err := model.ParseMoney(row[1])
			if err != nil {
				return fmt.Errorf("failed to parse credit amount '%s': %w", row[1], err)
			}

			debit, err := model.ParseMoney(row[2])
			if err != nil {
				return fmt.Errorf("failed to parse debit amount '%s': %w", row[2], err)
			}

			description := row[3]
			id := row[4]

			valueDate, err := time.Parse("02/01/2006", row[7])
			if err != nil {
				return fmt.Errorf("failed to parse value date '%s': %w", row[7], err)
			}

			var sourceAccountID, targetAccountID string
			var amount model.Money
			if credit > 0 {
				sourceAccountID = "external"
				targetAccountID = fmt.Sprintf("04-%s-%s", branch, account)
				amount = credit
			} else if debit > 0 {
				sourceAccountID = fmt.Sprintf("04-%s-%s", branch, account)
				targetAccountID = "external"
				amount = debit
			}

			inputs = append(inputs, &model.NewTransaction{
				Date:            valueDate,
				TargetAccountID: targetAccountID,
				SourceAccountID: sourceAccountID,
				ReferenceID:     id,
				Amount:          amount,
				Description:     description,
			})

			//if inputs != "" {
			//	inputs += ", "
			//}
			//inputs += fmt.Sprintf(
			//	`{date: "%s", targetAccountID: "%s", sourceAccountID: "%s", referenceID: "%s", amount: "%s", description: "%s"}`,
			//	valueDate.Format("2006-01-02"),
			//	targetAccountID,
			//	sourceAccountID,
			//	id,
			//	amount,
			//	description,
			//)
		}
	}

	//query := `{ mutation { createTransactions(inputs: [` + inputs + `]) } }`
	_, err = r.Mutation().CreateTransactions(ctx, inputs)
	if err != nil {
		return fmt.Errorf("failed to create transactions: %w", err)
	}
	return nil
}

func (r *Resolver) TransactionsXLSXProcessor() *transactionsXLSXProcessor {
	return &transactionsXLSXProcessor{r}
}

type transactionsXLSXProcessor struct{ *Resolver }
