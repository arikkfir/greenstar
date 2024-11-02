UPDATE tenants
SET updated_at   = now(),
    display_name = $2
WHERE id = $1
