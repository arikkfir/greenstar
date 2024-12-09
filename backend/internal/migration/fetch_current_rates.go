package migration

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	//go:embed insert_exchange_rate.sql
	insertUpdatedExchangeRateSQL string
)

func (m *exchangeRatesManagerImpl) UpdateExchangeRatesForToday(ctx context.Context) error {
	ctx, span := observability.Trace(ctx, trace.SpanKindServer)
	defer span.End()

	slog.DebugContext(ctx, "Updating today's exchange rates")

	// Fetch currencies
	rows, err := m.pool.Query(ctx, "SELECT code FROM currencies ORDER BY code")
	if err != nil {
		return fmt.Errorf("failed querying currencies: %w", err)
	}
	defer rows.Close()

	var currencyCodes []string
	for rows.Next() {
		var currencyCode string
		if err := rows.Scan(&currencyCode); err != nil {
			return fmt.Errorf("failed scanning currency code: %w", err)
		}
		currencyCodes = append(currencyCodes, currencyCode)
	}

	// Fetch exchange rates for each currency against all other currencies
	for _, currencyCode := range currencyCodes {
		if err := m.fetchTodayExchangeRatesForCurrency(ctx, currencyCode, currencyCodes); err != nil {
			return fmt.Errorf("failed fetching latest exchange rates for currency '%s': %w", currencyCode, err)
		}
	}

	return nil
}

func (m *exchangeRatesManagerImpl) fetchTodayExchangeRatesForCurrency(ctx context.Context, currencyCode string, currencies []string) error {
	ctx, span := observability.Trace(ctx, trace.SpanKindServer)
	defer span.End()

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed creating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	m.tokens.AcquireToken()

	query := url.Values{}
	query.Set("apikey", m.currencyAPIKey)
	query.Set("base_currency", currencyCode)
	query.Set("currencies", strings.Join(currencies, ","))
	query.Set("type", "fiat")
	response, err := http.Get("https://api.currencyapi.com/v3/latest?" + query.Encode())
	if err != nil {
		return fmt.Errorf("failed sending request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		bytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("request failed: %w\n%s", fmt.Errorf("unexpected status code %s (%d)", response.Status, response.StatusCode), string(bytes))
	}

	type exchangeRate struct {
		Code  string          `json:"code"`
		Value decimal.Decimal `json:"value"`
	}

	type latestExchangeRatesResponse struct {
		Data map[string]exchangeRate `json:"data"`
	}

	decoder := json.NewDecoder(response.Body)
	latestExchangeRates := latestExchangeRatesResponse{}
	if err := decoder.Decode(&latestExchangeRates); err != nil {
		return fmt.Errorf("failed decoding response: %w", err)
	}

	today := time.Now().UTC()
	for _, r := range latestExchangeRates.Data {
		if _, err := tx.Exec(ctx, insertUpdatedExchangeRateSQL, today, currencyCode, r.Code, r.Value, false, nil); err != nil {
			return fmt.Errorf("failed storing rate of %v from '%s' to '%s': %w", today, currencyCode, r.Code, err)
		}
		if _, err := tx.Exec(ctx, insertUpdatedExchangeRateSQL, today, currencyCode, r.Code, lang.DecimalOne.Div(r.Value), false, nil); err != nil {
			return fmt.Errorf("failed storing rate of %v from '%s' to '%s': %w", today, currencyCode, r.Code, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}
