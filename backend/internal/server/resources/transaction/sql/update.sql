UPDATE ONLY transactions AS t
SET updated_at        = now(),
    amount            = $1,
    currency          = $2,
    date              = $3,
    description       = $4,
    reference_id      = $5,
    source_account_id = $6,
    target_account_id = $7
FROM accounts AS sa,
     accounts AS ta
WHERE t.source_account_id = sa.id
  AND t.target_account_id = ta.id
  AND sa.tenant_id = $8
  AND ta.tenant_id = $8
  AND t.id = $9
