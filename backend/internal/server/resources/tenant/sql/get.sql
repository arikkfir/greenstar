SELECT id, created_at, updated_at, display_name
FROM tenants
WHERE id = $1
