UPDATE accounts
SET updated_at   = now(),
    display_name = $3,
    icon         = $4,
    parent_id    = $5
WHERE tenant_id = $1
  AND id = $2
