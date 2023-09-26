UPDATE transactions
SET updated_at = now()
WHERE id = $1
