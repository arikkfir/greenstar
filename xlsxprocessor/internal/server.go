package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/common"
	txPkg "github.com/arikkfir/greenstar/txprocessor/pkg"
	"github.com/arikkfir/greenstar/xlsxprocessor/pkg"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
	"regexp"
	"time"
)

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

type Server struct {
	Config Config
	Redis  *redis.Client
}

func NewServer(config Config, Redis *redis.Client) *Server {
	return &Server{
		Config: config,
		Redis:  Redis,
	}
}

func (s *Server) Run(ctx context.Context) error {
	logger := *zerolog.Ctx(ctx)
	subscriber := s.Redis.Subscribe(ctx, pkg.InputChannelName)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to receive message from Redis")
		}

		req := pkg.ProcessXLSXRequest{}
		if err := json.Unmarshal([]byte(msg.Payload), &req); err != nil {
			logger.Fatal().Err(err).Msg("Failed to unmarshal message")
		}

		if err := s.processXLSXFile(ctx, &req); err != nil {
			logger.Fatal().Err(err).Msg("Failed to process conversion request")
		}
	}
}

func (s *Server) processXLSXFile(ctx context.Context, req *pkg.ProcessXLSXRequest) error {
	f, err := excelize.OpenReader(bytes.NewReader(req.Data))
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}

	if f.SheetCount != 2 {
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

	if rows, err := f.GetRows(transactionsSheet.Name, excelize.Options{}); err != nil {
		return fmt.Errorf("failed to get rows from sheet '%s': %w", transactionsSheet.Name, err)
	} else if len(rows) < 7 {
		// no transactions
		return nil
	} else {
		for _, row := range rows[6:] {
			if len(row) == 10 {
				credit, err := common.ParseCurrency(row[1])
				if err != nil {
					return fmt.Errorf("failed to parse credit amount '%s': %w", row[1], err)
				}

				debit, err := common.ParseCurrency(row[2])
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
				var amount common.Currency
				if credit > 0 {
					sourceAccountID = "external"
					targetAccountID = fmt.Sprintf("04-%s-%s", branch, account)
					amount = credit
				} else if debit > 0 {
					sourceAccountID = fmt.Sprintf("04-%s-%s", branch, account)
					targetAccountID = "external"
					amount = debit
				}
				tx := txPkg.Transaction{
					Date:            valueDate.Format("2006-01-02"),
					SourceAccountID: sourceAccountID,
					TargetAccountID: targetAccountID,
					Amount:          amount,
					Description:     description,
					ReferenceID:     id,
				}
				if payload, err := json.Marshal(&tx); err != nil {
					return fmt.Errorf("failed to marshal transaction: %w", err)
				} else if err := s.Redis.Publish(ctx, pkg.TargetChannelName, payload).Err(); err != nil {
					return fmt.Errorf("failed to publish transaction to Redis: %w", err)
				}
			}
		}
		return nil
	}
}
