INSERT INTO transactions (amount, currency, date, description, reference_id, source_account_id, target_account_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
