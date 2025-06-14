import { Paper } from "@mui/material"
import "./DashboardPage.scss"
import { AccountNode, AccountsTree } from "../../components/AccountsTree.tsx"
import { useCallback, useState } from "react"
import { AccountBalanceChart } from "./AccountBalanceChart.tsx"

export function DashboardPage() {
    const [ selectedAccountIDs, setSelectedAccountIDs ] = useState<string[]>([])

    const handleAccountSelectionChange = useCallback(
        async (accounts: AccountNode[]) => setSelectedAccountIDs(accounts.map(a => a.id)),
        [],
    )

    return (
        <main className="dashboard-page">
            <title>Home - GreenSTAR</title>
            <Paper className="accounts-tree-container" elevation={3}>
                <AccountsTree<true> onAccountSelectionChange={handleAccountSelectionChange}
                                    checkboxSelection multiSelect />
            </Paper>

            <AccountBalanceChart accountIDs={selectedAccountIDs} />

            {/*<Paper className="expenses-chart-container" elevation={3}>*/}
            {/*    <ExpensesPieChart endDate={lastTransactionDate} />*/}
            {/*</Paper>*/}

            {/*<Paper className="top-accounts-container" elevation={3}>*/}
            {/*    <TopAccountsGrid endDate={lastTransactionDate} />*/}
            {/*</Paper>*/}
        </main>
    )
}
