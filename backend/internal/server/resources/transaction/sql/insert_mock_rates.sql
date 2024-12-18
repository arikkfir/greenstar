INSERT INTO rates (date, source_currency_code, target_currency_code, rate, mock)
SELECT date($1), $2, currencies.code, 1, true
FROM currencies
UNION
SELECT date($1), currencies.code, $2, 1, true
FROM currencies
ON CONFLICT ON CONSTRAINT rates_pkey DO NOTHING
