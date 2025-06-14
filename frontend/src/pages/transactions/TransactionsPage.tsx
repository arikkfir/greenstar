import "./Layout.scss"
import { useCallback, useState } from "react"
import { Paper } from "@mui/material"
import { AccountNode, AccountsTree } from "../../components/AccountsTree.tsx"
import { TransactionsTable } from "./TransactionsTable.tsx"

export function TransactionsPage() {

    // Selected account nodes state
    const [ selectedAccountIDs, setSelectedAccountIDs ] = useState<string | undefined>()

    // Callback to handle account selection from the AccountsTree
    const handleAccountSelectionChange = useCallback(
        (accounts: AccountNode[]): void => setSelectedAccountIDs(accounts.length ? accounts[0].id : undefined),
        [ setSelectedAccountIDs ],
    )

    return (
        <main className="transactions-page">
            <title>Transactions - GreenSTAR</title>

            <Paper className="accounts-tree-container" elevation={3}>
                <AccountsTree<false> showBalance onAccountSelectionChange={handleAccountSelectionChange} />
            </Paper>

            <Paper className="transactions-container" elevation={3}>
                <TransactionsTable involvingAccountID={selectedAccountIDs} />
            </Paper>

        </main>
    )
}
