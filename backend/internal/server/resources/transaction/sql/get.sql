SELECT t.id,
       t.created_at,
       t.updated_at,
       t.amount,
       t.currency,
       t.amount * rates.rate AS converted_amount,
       t.date,
       t.description,
       t.reference_id,
       t.source_account_id,
       t.target_account_id
FROM transactions AS t
         JOIN accounts AS sa ON t.source_account_id = sa.id
         JOIN accounts AS ta ON t.target_account_id = ta.id
         JOIN rates ON date(t.date AT TIME ZONE 'UTC') = rates.date AND rates.source_currency_code = t.currency AND
                       rates.target_currency_code = $3
WHERE t.id = $1
  AND sa.tenant_id = $2
  AND ta.tenant_id = $2
