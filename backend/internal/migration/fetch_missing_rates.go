package migration

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/util/lang"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

var (
	//go:embed fetch_missing_rates.sql
	fetchMissingRatesSQL string

	//go:embed insert_exchange_rate.sql
	insertMissingExchangeRateSQL string
)

type missingRate struct {
	Date                time.Time
	BaseCurrencyCode    string
	TargetCurrencyCodes string
}

func (m *exchangeRatesManagerImpl) PopulateMissingExchangeRates(ctx context.Context) error {
	ctx, span := observability.Trace(ctx, trace.SpanKindServer)
	defer span.End()

	slog.Default().InfoContext(ctx, "Searching for missing historical exchange rates based on existing transactions")

	rows, err := m.pool.Query(ctx, fetchMissingRatesSQL)
	if err != nil {
		return fmt.Errorf("failed fetching missing rates: %w", err)
	}
	defer rows.Close()

	var missingRates []missingRate
	for rows.Next() {
		var r missingRate
		if err := rows.Scan(&r.Date, &r.BaseCurrencyCode, &r.TargetCurrencyCodes); err != nil {
			return fmt.Errorf("failed scanning missing rate: %w", err)
		}
		missingRates = append(missingRates, r)
	}
	if len(missingRates) == 0 {
		return nil
	}

	slog.Default().DebugContext(ctx, "Found missing exchange rates - obtaining current quotes", "missingCount", len(missingRates))
	for _, r := range missingRates {
		if err := m.fetchExchangeRates(ctx, r); err != nil {
			return fmt.Errorf("failed processing missing rate '%+v': %w", r, err)
		}
	}

	return nil
}

func (m *exchangeRatesManagerImpl) fetchExchangeRates(ctx context.Context, r missingRate) error {
	ctx, span := observability.Trace(ctx, trace.SpanKindServer)
	defer span.End()

	m.tokens.AcquireToken()

	query := url.Values{}
	query.Set("apikey", m.currencyAPIKey)
	query.Set("base_currency", r.BaseCurrencyCode)
	query.Set("currencies", r.TargetCurrencyCodes)
	query.Set("date", r.Date.Format("2006-01-02"))
	query.Set("type", "fiat")
	response, err := http.Get("https://api.currencyapi.com/v3/historical?" + query.Encode())
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

	type historicalExchangeRatesResponse struct {
		Data map[string]exchangeRate `json:"data"`
	}

	decoder := json.NewDecoder(response.Body)
	historicalExchangeRates := historicalExchangeRatesResponse{}
	if err := decoder.Decode(&historicalExchangeRates); err != nil {
		return fmt.Errorf("failed decoding response: %w", err)
	}

	for _, rate := range historicalExchangeRates.Data {
		if _, err := m.pool.Exec(ctx, insertMissingExchangeRateSQL, r.Date, r.BaseCurrencyCode, rate.Code, rate.Value, false, nil); err != nil {
			return fmt.Errorf("failed storing rate '%+v': %w", r, err)
		}
		if _, err := m.pool.Exec(ctx, insertMissingExchangeRateSQL, r.Date, rate.Code, r.BaseCurrencyCode, lang.DecimalOne.Div(rate.Value), false, nil); err != nil {
			return fmt.Errorf("failed storing inverse rate '%+v': %w", r, err)
		}
	}

	return nil
}
