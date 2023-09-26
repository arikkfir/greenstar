package migration

import (
	"context"
	"github.com/arikkfir/greenstar/backend/internal/util/token"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// TODO: obtain a lease to ensure only one replica initializes exchange rates
// TODO: add a background process that uploads current rates table to the remote bucket, so next historical rates download is "fuller"

type HistoricalExchangeRatesPeriod string

const (
	Full   HistoricalExchangeRatesPeriod = "full"
	Latest HistoricalExchangeRatesPeriod = "latest"

	historicalExchangeRatesBucketName    = "arikkfir-greenstar"
	historicalExchangeRatesTarObjectName = "rates.csv"
)

type ExchangeRatesManager interface {
	PopulateCurrencies(ctx context.Context) error
	PopulateHistoricalRates(context.Context, HistoricalExchangeRatesPeriod) error
	PopulateMissingExchangeRates(context.Context) error
}

type exchangeRatesManagerImpl struct {
	currencyAPIKey string
	pool           *pgxpool.Pool
	tokens         *token.Bucket
}

func NewExchangeRatesManager(currencyAPIKey string, pool *pgxpool.Pool) ExchangeRatesManager {
	return &exchangeRatesManagerImpl{
		currencyAPIKey: currencyAPIKey,
		pool:           pool,
		tokens:         token.NewBucket(10, time.Minute/time.Duration(10)), // 1 token every 6sec; i.e. 60/10
	}
}
