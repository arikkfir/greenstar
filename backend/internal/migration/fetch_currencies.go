package migration

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/backend/internal/util/observability"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	SupportedCurrencies = []string{
		"USD",
		"EUR",
		"JPY",
		"GBP",
		"CNY",
		"CAD",
		"AUD",
		"CHF",
		"HKD",
		"SEK",
		"SGD",
		"PLN",
		"NOK",
		"DKK",
		"NZD",
		"MXN",
		"CZK",
		"ZAR",
		"HUF",
		"THB",
		"SAR",
		"ILS",
		"RUB",
		"TRY",
	}
)

var (
	//go:embed insert_currency.sql
	insertCurrencySQL string
)

// PopulateCurrencies fetches and updates the system with the latest currency information
// for currencies that are missing or stale. It connects to an external currency API,
// retrieves data, and updates the local database with this information.
func (m *exchangeRatesManagerImpl) PopulateCurrencies(ctx context.Context) error {
	ctx, span := observability.Trace(ctx, trace.SpanKindServer)
	defer span.End()

	slog.InfoContext(ctx, "Updating system currencies")

	staleCurrencies, err := m.fetchMissingOrStaleCurrencies(ctx)
	if err != nil {
		return fmt.Errorf("failed fetching stale currencies: %w", err)
	} else if len(staleCurrencies) == 0 {
		return nil
	}

	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed creating transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	m.tokens.AcquireToken()

	query := url.Values{}
	query.Set("apikey", m.currencyAPIKey)
	query.Set("currencies", strings.Join(staleCurrencies, ","))
	response, err := http.Get("https://api.currencyapi.com/v3/currencies?" + query.Encode())
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		bytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("currencies request failed: %w\n%s", fmt.Errorf("unexpected status code %s (%d)", response.Status, response.StatusCode), string(bytes))
	}

	type currencyInfo struct {
		Symbol        string   `json:"symbol"`
		SymbolNative  string   `json:"symbol_native"`
		Name          string   `json:"name"`
		DecimalDigits int      `json:"decimal_digits"`
		Rounding      int      `json:"rounding"`
		Code          string   `json:"code"`
		NamePlural    string   `json:"name_plural"`
		Type          string   `json:"type"`
		CountryCodes  []string `json:"countries"`
	}

	type currenciesListResponse struct {
		Data map[string]currencyInfo `json:"data"`
	}

	decoder := json.NewDecoder(response.Body)
	currenciesList := currenciesListResponse{}
	if err := decoder.Decode(&currenciesList); err != nil {
		return fmt.Errorf("failed decoding response: %w", err)
	}

	var currencies []string
	for _, c := range currenciesList.Data {
		symbol := c.Symbol
		symbolNative := c.SymbolNative
		name := c.Name
		decimalDigits := c.DecimalDigits
		rounding := c.Rounding
		code := c.Code
		namePlural := c.NamePlural
		currencyType := c.Type
		countryCodes := c.CountryCodes
		if _, err := tx.Exec(ctx, insertCurrencySQL, symbol, symbolNative, name, decimalDigits, rounding, code, namePlural, currencyType, countryCodes); err != nil {
			return fmt.Errorf("failed storing currency %s: %w", code, err)
		}
		currencies = append(currencies, code)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	return nil
}

// fetchMissingOrStaleCurrencies identifies currencies that are either not present or need updates
// based on their last update time. It returns a list of those currency codes, or an error if one occurs.
func (m *exchangeRatesManagerImpl) fetchMissingOrStaleCurrencies(ctx context.Context) ([]string, error) {
	ctx, span := observability.Trace(ctx, trace.SpanKindServer)
	defer span.End()

	// Start with all supported currencies, assuming all are stale (to ensure we fetch missing ones...)
	currencies := make(map[string]time.Time)
	for _, supportedCurrency := range SupportedCurrencies {
		currencies[supportedCurrency] = time.Now().Add(-1 * 24 * 30 * time.Hour)
	}

	// Fetch all currencies from the database
	rows, err := m.pool.Query(ctx, "SELECT code, updated_at FROM currencies")
	if err != nil {
		return nil, fmt.Errorf("failed querying currencies: %w", err)
	}
	defer rows.Close()

	// Update the stale currencies map with the correct update date from the fetched currencies
	for rows.Next() {
		var currencyCode string
		var updatedAt time.Time
		if err := rows.Scan(&currencyCode, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed scanning currency code: %w", err)
		}
		currencies[currencyCode] = updatedAt
	}

	// Collect the stale ones
	var staleCurrencies []string
	for code, updatedAt := range currencies {
		if time.Since(updatedAt) > 24*time.Hour {
			staleCurrencies = append(staleCurrencies, code)
		}
	}

	return staleCurrencies, nil
}
