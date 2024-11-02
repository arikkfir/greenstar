WITH RECURSIVE
    AccountHierarchy AS (SELECT id AS root_account_id,
                                id AS account_id
                         FROM accounts
                         UNION ALL
                         SELECT ah.root_account_id,
                                a.id AS account_id
                         FROM AccountHierarchy ah
                                  JOIN accounts a ON a.parent_id = ah.account_id),
    AccountSums AS (SELECT ah.root_account_id,
                           SUM(
                                   CASE
                                       WHEN t.target_account_id = ah.account_id THEN t.amount * rates.rate
                                       ELSE 0
                                       END
                           ) AS total_in,
                           SUM(
                                   CASE
                                       WHEN t.source_account_id = ah.account_id THEN t.amount * rates.rate
                                       ELSE 0
                                       END
                           ) AS total_out
                    FROM AccountHierarchy ah
                             LEFT JOIN transactions t
                                       ON ah.account_id = t.source_account_id OR ah.account_id = t.target_account_id
                             LEFT JOIN rates
                                       ON date(t.date AT TIME ZONE 'UTC') = rates.date
                                           AND rates.source_currency_code = t.currency
                                           AND rates.target_currency_code = $2
                    GROUP BY ah.root_account_id)
SELECT a.id,
       a.created_at,
       a.updated_at,
       a.display_name,
       a.icon,
       a.parent_id,
       COALESCE(account_balance.total_in, 0)                                          AS incoming,
       COALESCE(account_balance.total_out, 0)                                         AS outgoing,
       COALESCE(account_balance.total_in, 0) - COALESCE(account_balance.total_out, 0) AS balance
FROM accounts a
         LEFT JOIN AccountSums account_balance ON a.id = account_balance.root_account_id
WHERE a.tenant_id = $1
