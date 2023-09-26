UPDATE tenants
SET updated_at = now()
WHERE id = $1
