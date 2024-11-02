DELETE
FROM ONLY transactions
    USING accounts AS sa, accounts AS ta
WHERE transactions.source_account_id = sa.id
  AND transactions.target_account_id = ta.id
  AND sa.tenant_id = $1
  AND ta.tenant_id = $1
