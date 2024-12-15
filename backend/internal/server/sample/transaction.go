package sample

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/server/resources/transaction"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/strings"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"
	"text/template"
	"time"
)

type transactionSpec struct {
	TargetAccountID string          `yaml:"targetAccountID"`
	Date            time.Time       `yaml:"date"`
	Amount          decimal.Decimal `yaml:"amount"`
	Currency        string          `yaml:"currency"`
	Description     string          `yaml:"description"`
	Repeat          repeatSpec      `yaml:"repeat"`
}

func (t *transactionSpec) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		TargetAccountID string      `yaml:"targetAccountID"`
		Date            string      `yaml:"date"`
		Amount          interface{} `yaml:"amount"`
		Currency        string      `yaml:"currency"`
		Description     string      `yaml:"description"`
		Repeat          repeatSpec  `yaml:"repeat"`
	}

	if err := value.Decode(&raw); err != nil {
		return err
	}

	if d, err := timeFromString("date", raw.Date); err != nil {
		return fmt.Errorf("failed decoding YAML transaction date: %w", err)
	} else {
		t.Date = d
	}

	if amount, err := decimalFromValue("amount", raw.Amount); err != nil {
		return fmt.Errorf("failed decoding YAML transaction amount: %w", err)
	} else {
		t.Amount = amount
	}

	t.TargetAccountID = raw.TargetAccountID
	t.Currency = raw.Currency
	t.Description = raw.Description
	t.Repeat = raw.Repeat
	return nil
}

func (t *transactionSpec) apply(ctx context.Context, ids map[string]string, h transaction.Handler, defaultCurrency, sourceAccountID string) error {
	sourceAccountUUID := ids[sourceAccountID]
	if sourceAccountUUID == "" {
		return fmt.Errorf("could not find UUID for source account ID '%s'", sourceAccountID)
	}

	targetAccountUUID := ids[t.TargetAccountID]
	if targetAccountUUID == "" {
		return fmt.Errorf("could not find UUID for target account ID '%s'", t.TargetAccountID)
	}

	txDate := t.Date
	txAmount := t.Amount
	txCurrency := cmp.Or(t.Currency, defaultCurrency)
	if _, err := h.Create(ctx, transaction.CreateRequest{
		Date:            txDate,
		ReferenceID:     strings.RandomHash(7),
		Amount:          txAmount,
		Currency:        txCurrency,
		Description:     lang.PtrOf(t.Description),
		SourceAccountID: sourceAccountUUID,
		TargetAccountID: targetAccountUUID,
	}); err != nil {
		return fmt.Errorf("failed creating transaction: %w", err)
	}
	if err := t.Repeat.apply(ctx, ids, h, txCurrency, t.Description, sourceAccountID, t.TargetAccountID, txDate, txAmount); err != nil {
		return fmt.Errorf("failed repeating transaction: %w", err)
	}
	return nil
}

type repeatSpec struct {
	Count          int                `yaml:"count"`
	DateTemplate   *template.Template `yaml:"dateTemplate"`
	AmountTemplate *template.Template `yaml:"amountTemplate"`
}

func (r *repeatSpec) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Count          int    `yaml:"count"`
		DateTemplate   string `yaml:"dateTemplate"`
		AmountTemplate string `yaml:"amountTemplate"`
	}

	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding YAML repeat spec: %w", err)
	}

	if raw.DateTemplate != "" {
		if tmpl, err := newTemplate("date", raw.DateTemplate); err != nil {
			return fmt.Errorf("failed decoding transaction repeat-spec's YAML date template: %w", err)
		} else {
			r.DateTemplate = tmpl
		}
	}
	if raw.AmountTemplate != "" {
		if tmpl, err := newTemplate("amount", raw.AmountTemplate); err != nil {
			return fmt.Errorf("failed decoding transaction repeat-spec's YAML amount template: %w", err)
		} else {
			r.AmountTemplate = tmpl
		}
	}

	r.Count = raw.Count
	return nil
}

func (r *repeatSpec) apply(ctx context.Context, ids map[string]string, h transaction.Handler, currency, description, sourceAccountID, targetAccountID string, baseDate time.Time, baseAmount decimal.Decimal) error {
	sourceAccountUUID := ids[sourceAccountID]
	if sourceAccountUUID == "" {
		return fmt.Errorf("could not find UUID for source account ID '%s'", sourceAccountID)
	}

	targetAccountUUID := ids[targetAccountID]
	if targetAccountUUID == "" {
		return fmt.Errorf("could not find UUID for target account ID '%s'", targetAccountID)
	}

	var res bytes.Buffer
	for i := 0; i < r.Count; i++ {
		if r.DateTemplate != nil {
			res.Reset()
			if err := r.DateTemplate.Execute(&res, baseDate); err != nil {
				return fmt.Errorf("failed modifying transaction date: %w", err)
			} else if d, err := time.Parse(time.RFC3339, res.String()); err != nil {
				return fmt.Errorf("failed decoding modified transaction date '%s' into RFC3339: %w", res.String(), err)
			} else {
				baseDate = d
			}
		}
		if r.AmountTemplate != nil {
			res.Reset()
			if err := r.AmountTemplate.Execute(&res, baseAmount); err != nil {
				return fmt.Errorf("failed modifying transaction amount: %w", err)
			} else if d, err := decimal.NewFromString(res.String()); err != nil {
				return fmt.Errorf("failed decoding modified transaction amount '%s': %w", res.String(), err)
			} else {
				baseAmount = d
			}
		}
		if _, err := h.Create(ctx, transaction.CreateRequest{
			Date:            baseDate,
			ReferenceID:     strings.RandomHash(7),
			Amount:          baseAmount,
			Currency:        currency,
			Description:     lang.PtrOf(description),
			SourceAccountID: sourceAccountUUID,
			TargetAccountID: targetAccountUUID,
		}); err != nil {
			return fmt.Errorf("failed creating repeated transaction: %w", err)
		}
	}
	return nil
}
