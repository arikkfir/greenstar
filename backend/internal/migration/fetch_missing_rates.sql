SELECT date(t.date AT TIME ZONE 'UTC') AS date,
       t.currency                      AS base_currency,
       string_agg(c.code, ',')         AS target_currencies
FROM transactions AS t
         JOIN currencies AS c ON t.currency != c.code
         LEFT JOIN rates
                   ON date(t.date AT TIME ZONE 'UTC') = rates.date
                       AND rates.source_currency_code = t.currency AND rates.target_currency_code = c.code
WHERE (rates.date IS NULL OR rates.mock)
  AND t.date::date < now()::date
GROUP BY date(t.date AT TIME ZONE 'UTC'), t.currency
ORDER BY date(t.date AT TIME ZONE 'UTC'), t.currency;
