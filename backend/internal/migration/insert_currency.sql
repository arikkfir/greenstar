INSERT INTO currencies (symbol,
                        symbol_native,
                        name,
                        decimal_digits,
                        rounding,
                        code,
                        name_plural,
                        type,
                        country_codes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT ON CONSTRAINT currencies_pkey
    DO UPDATE SET symbol         = $1,
                  symbol_native  = $2,
                  name           = $3,
                  decimal_digits = $4,
                  rounding       = $5,
                  name_plural    = $7,
                  type           = $8,
                  country_codes  = $9,
                  updated_at     = now();
