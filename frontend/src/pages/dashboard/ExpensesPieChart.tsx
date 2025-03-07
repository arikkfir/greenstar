import { useCallback, useContext, useEffect, useMemo, useState } from "react"
import { Card, CardContent, CardHeader } from "@mui/material"
import { PieChart } from "@mui/x-charts/PieChart"
import { AccountData, AccountSelector } from "./AccountSelector"
import { useLazyQuery } from "@apollo/client"
import { gql } from "../../graphql"
import { useTenantID } from "../../hooks/tenant.ts"
import { SelectedCurrencyContext } from "../../contexts/SelectedCurrency.ts"
import { DateTime } from "luxon"
import { Account } from "../../graphql/graphql.ts"
import { useCurrencyFormatter } from "../../hooks/locale.ts"

const SubAccounts = gql(`
    query SubAccountsQuery($tenantID: ID!, $accountID: ID!, $currency: String!, $until: DateTime) {
        tenant(id: $tenantID) {
            id
            account(id: $accountID) {
                id
                children {
                    id
                    displayName
                    icon
                    balance(currency: $currency, until: $until)
                }
            }
        }
    }
`)

interface ExpenseData {
    id: string;
    label: string;
    value: number;
    color?: string;
}

interface ExpensesPieChartProps {
    endDate?: DateTime;
}

export function ExpensesPieChart({ endDate }: ExpensesPieChartProps) {
    const tenantID                                = useTenantID()
    const selectedCurrency                        = useContext(SelectedCurrencyContext)
    const currencyFormatter                       = useCurrencyFormatter()
    const [ selectedAccount, setSelectedAccount ] = useState<AccountData | null>(null)
    const [ expensesData, setExpensesData ]       = useState<ExpenseData[]>([])
    const [ getSubAccounts, { loading } ]         = useLazyQuery(SubAccounts)

    // Fetch child accounts data for the selected account
    const fetchExpensesData = useCallback(
        async () => {
            if (!selectedAccount || !selectedCurrency.currency) {
                setExpensesData([])
                return
            }

            // Fetch
            const result = await getSubAccounts({
                variables: {
                    tenantID,
                    accountID: selectedAccount.id,
                    currency: selectedCurrency.currency.code,
                    until: endDate,
                },
            })

            if (result.data?.tenant?.account?.children) {
                // Filter child accounts with negative balances (expenses)
                const expenses = result.data.tenant.account.children
                                       .filter((account: { balance: Account["balance"] }) => account.balance < 0)
                                       .map((account: {
                                           id: Account["id"],
                                           displayName: Account["displayName"],
                                           balance: Account["balance"]
                                       }) => ({
                                           id: account.id,
                                           label: account.displayName,
                                           value: Math.abs(account.balance), // Use absolute value for the pie chart
                                       }))
                setExpensesData(expenses)
            } else {
                setExpensesData([])
            }
        },
        [ selectedAccount, selectedCurrency.currency, tenantID, getSubAccounts, endDate ],
    )

    // Fetch expenses data when the selected account changes
    useEffect(
        () => {
            fetchExpensesData().then()
        },
        [ selectedAccount, fetchExpensesData ],
    )

    // Prepare data for the chart
    const chartData = useMemo(
        () => {
            // If there are no expenses, return an empty array
            if (expensesData.length === 0) {
                return []
            }

            // If there are more than 5 expenses, group the smallest ones into "Other"
            if (expensesData.length > 5) {
                // Sort by value (largest first)
                const sorted = [ ...expensesData ].sort((a, b) => b.value - a.value)

                // Take top 4 expenses
                const topExpenses = sorted.slice(0, 4)

                // Sum the rest into "Other"
                const otherExpenses = sorted.slice(4)
                const otherValue    = otherExpenses.reduce((sum, expense) => sum + expense.value, 0)

                return [
                    ...topExpenses,
                    {
                        id: "other",
                        label: "Other",
                        value: otherValue,
                    },
                ]
            }

            return expensesData
        },
        [ expensesData ],
    )

    return (
        <Card>
            <CardHeader title="Expenses Distribution"
                        action={
                            <AccountSelector storageKey="dashboard-expenses-chart-account"
                                             label="Select Account"
                                             onChange={setSelectedAccount} />
                        } />
            <CardContent>
                {chartData.length > 0 ? (
                    <PieChart
                        series={[
                            {
                                data: chartData,
                                highlightScope: { highlight: "item" },
                                arcLabel: (value) => {
                                    const currencyCode = selectedCurrency?.currency?.code || "USD"
                                    return currencyFormatter(value.value, currencyCode, 0, 2)
                                },
                            },
                        ]}
                        height={300}
                        margin={{ top: 0, bottom: 0, left: 0, right: 0 }}
                        slotProps={{
                            legend: {
                                direction: "vertical",
                                position: { vertical: "middle", horizontal: "end" },
                            },
                        }}
                    />
                ) : (
                    <div style={{ height: 300, display: "flex", alignItems: "center", justifyContent: "center" }}>
                        {loading ? "Loading expenses data..." : "Select an account to view expenses distribution"}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
