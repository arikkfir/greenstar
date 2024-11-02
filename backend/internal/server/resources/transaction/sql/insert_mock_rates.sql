INSERT INTO rates (date, source_currency_code, target_currency_code, rate, mock)
VALUES ($1, $2, $3, 1, true)
ON CONFLICT ON CONSTRAINT rates_pkey DO NOTHING
