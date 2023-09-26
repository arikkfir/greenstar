WITH RECURSIVE
    source_accounts AS (SELECT id, parent_id
                        FROM accounts
                        WHERE tenant_id = $3
                          AND id = $1
                        UNION ALL
                        SELECT a.id, a.parent_id
                        FROM accounts a
                                 INNER JOIN source_accounts AS wsa ON a.parent_id = wsa.id
                        WHERE tenant_id = $3),
    target_accounts AS (SELECT id, parent_id
                        FROM accounts
                        WHERE tenant_id = $3
                          AND id = $2
                        UNION ALL
                        SELECT a.id, a.parent_id
                        FROM accounts a
                                 INNER JOIN target_accounts AS wta ON a.parent_id = wta.id
                        WHERE tenant_id = $3)
SELECT t.id,
       t.created_at,
       t.updated_at,
       t.date,
       t.reference_id,
       t.amount,
       t.currency,
       t.amount * rates.rate AS converted_amount,
       t.description,
       t.source_account_id,
       t.target_account_id
FROM transactions t
         JOIN accounts sa ON t.source_account_id = sa.id
         JOIN accounts ta ON t.target_account_id = ta.id
         JOIN rates ON date(t.date AT TIME ZONE 'UTC') = rates.date AND rates.source_currency_code = t.currency AND
                       rates.target_currency_code = $4
WHERE sa.tenant_id = $3
  AND ta.tenant_id = $3
