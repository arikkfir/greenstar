DELETE
FROM accounts
WHERE tenant_id = $1
  AND id = $2
