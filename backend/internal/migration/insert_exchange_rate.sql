INSERT INTO rates (date, source_currency_code, target_currency_code, rate, mock, line_number)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ON CONSTRAINT rates_pkey
    DO UPDATE SET rate        = $4,
                  mock        = $5,
                  line_number = $6,
                  updated_at  = now();