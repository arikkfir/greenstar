import { useCallback, useContext, useEffect, useMemo, useState } from "react"
import { LocaleContext } from "../providers/LocaleProvider.tsx"
import { Box, Paper } from "@mui/material"
import { AccountsTree } from "../components/tree/accounts/AccountsTree.tsx"
import { TransactionsTable } from "../components/table/transactions/TransactionsTable.tsx"
import {
    AccountID,
    AccountNode,
    BuildAccountsTree,
    FindAccountInTree,
    MapAccountsByID,
} from "../models/account-addons.ts"
import { useListAccounts } from "../hooks/account/list.ts"

export function TransactionsPage() {
    const locale = useContext(LocaleContext)
    const { request: fetchAccounts, operation: accountsStatus } = useListAccounts()
    const [selectedAccount, setSelectedAccount] = useState<AccountNode | undefined>()
    const accountsMap = useMemo(() => MapAccountsByID(accountsStatus.response?.items || []), [accountsStatus])
    const accountsTree = useMemo(() => BuildAccountsTree(accountsStatus.response?.items || []), [accountsStatus])

    useEffect(() => fetchAccounts({ currency: locale.currency }), [locale])

    const handleAccountSelectionChange = useCallback(
        (_: any, id: AccountID | null) => id && setSelectedAccount(FindAccountInTree(accountsTree, id)),
        [accountsTree, setSelectedAccount],
    )

    if (accountsStatus.error) {
        return <p>Error : {accountsStatus.error.message}</p>
    }

    return (
        <Box
            sx={{
                flexGrow: 1,
                flexShrink: 1,
                display: "flex",
                flexDirection: "row",
                justifyContent: "center",
                alignItems: "stretch",
                alignContent: "stretch",
                gap: "1rem",
                p: 2,
                overflow: "hidden",
            }}
        >
            <Paper sx={{ flexGrow: 0, flexShrink: 0, minWidth: "20rem", p: 1, overflow: "scroll" }}>
                <AccountsTree
                    accounts={accountsTree}
                    selectedItems={selectedAccount?.id || null}
                    onSelectedItemsChange={handleAccountSelectionChange}
                />
            </Paper>
            <Paper
                sx={{
                    flexGrow: 1,
                    flexShrink: 1,
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "stretch",
                    alignContent: "stretch",
                    overflow: "hidden",
                }}
            >
                <TransactionsTable
                    sx={{ flexGrow: 1, flexShrink: 1 }}
                    stateId="transactions"
                    sourceAccountId={selectedAccount?.id}
                    targetAccountId={selectedAccount?.id}
                    accounts={accountsMap}
                />
            </Paper>
        </Box>
    )
}
