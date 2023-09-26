INSERT INTO accounts (display_name, icon, parent_id, tenant_id)
VALUES ($1, $2, $3, $4)
RETURNING id
