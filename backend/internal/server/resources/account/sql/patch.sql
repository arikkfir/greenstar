UPDATE accounts
SET updated_at = now()
WHERE tenant_id = $1
  AND id = $2
