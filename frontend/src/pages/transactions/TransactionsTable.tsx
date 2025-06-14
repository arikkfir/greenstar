import { useLazyQuery } from "@apollo/client"
import { useMemo } from "react"
import {
    DataGridPremium,
    GridColDef,
    GridDataSource,
    GridGetRowsParams,
    GridGetRowsResponse,
} from "@mui/x-data-grid-premium"
import { SortDirection, TransactionsSortColumns, TransactionsSortColumnsInput } from "../../graphql/graphql.ts"
import { TransactionRow } from "./TransactionRow.tsx"
import { useTransactionColumns } from "./TransactionColumns.tsx"
import { useTenantID } from "../../hooks/tenant.ts"
import { DateTime } from "luxon"
import { Snackbar } from "@mui/material"
import { gql } from "../../graphql"

const ASC        = SortDirection.Asc
const DESC       = SortDirection.Desc
const TxSortCols = TransactionsSortColumns

const Transactions = gql(`
    query TransactionsQuery($tenantID: ID!, $involvingAccountID: ID, $offset: Int, $limit: Int, $sort: [TransactionsSortColumnsInput!], $until: DateTime) {
        tenant(id: $tenantID) {
            id
            transactions(involvingAccountID: $involvingAccountID, sort: $sort, offset: $offset, limit: $limit, until: $until) {
                rows {
                    id
                    createdAt
                    updatedAt
                    currency {
                        code
                    }
                    amount
                    date
                    description
                    referenceID
                    sourceAccount {
                        id
                        icon
                        displayName
                    }
                    targetAccount {
                        id
                        icon
                        displayName
                    }
                }
                totalCount
            }
        }
    }
`)

export interface TransactionsTableProps {
    involvingAccountID?: string
    until?: DateTime
}

export function TransactionsTable({ involvingAccountID, until }: TransactionsTableProps) {
    const tenantID = useTenantID()

    // Transactions grid columns (static)
    const transactionColumns: GridColDef<TransactionRow>[] = useTransactionColumns()

    // Transactions GraphQL query to fetch transactions; invoked by the grid data source on startup/account selection
    const [ fetchTransactions, fetchTransactionsResult ] = useLazyQuery(Transactions, {
        refetchWritePolicy: "overwrite",
        notifyOnNetworkStatusChange: true,
        variables: { tenantID },
    })

    const transactionsDataSource: GridDataSource = useMemo(
        (): GridDataSource => ({
            getRows: async (params: GridGetRowsParams): Promise<GridGetRowsResponse> => {
                if (params.filterModel.items.length) {
                    throw new Error("Filtering not implemented by server")
                } else if (params.filterModel.quickFilterValues?.length) {
                    throw new Error("Quick filter not implemented by server")
                } else if (params.groupKeys?.length) {
                    throw new Error("Grouped columns not implemented by server")
                } else if (params.aggregationModel?.length) {
                    throw new Error("Aggregation not implemented by server")
                }

                const sort: TransactionsSortColumnsInput[] = []
                for (let s of params.sortModel) {
                    switch (s.field) {
                        case "date":
                            sort.push({ col: TxSortCols.Date, direction: s.sort == "desc" ? DESC : ASC })
                            break
                        case "targetAccount":
                            sort.push({ col: TxSortCols.TargetAccountName, direction: s.sort == "desc" ? DESC : ASC })
                            break
                        case "sourceAccount":
                            sort.push({ col: TxSortCols.SourceAccountName, direction: s.sort == "desc" ? DESC : ASC })
                            break
                        case "description":
                            sort.push({ col: TxSortCols.Description, direction: s.sort == "desc" ? DESC : ASC })
                            break
                        case "amount":
                            sort.push({ col: TxSortCols.Amount, direction: s.sort == "desc" ? DESC : ASC })
                            break
                        case "referenceID":
                            sort.push({ col: TxSortCols.ReferenceId, direction: s.sort == "desc" ? DESC : ASC })
                            break
                    }
                }

                const result = await fetchTransactions({
                    variables: {
                        tenantID,
                        sort,
                        involvingAccountID,
                        until,
                        limit: params.end - (params.start as number) + 1,
                        offset: params.start as number,
                    },
                })

                if (result.error) {
                    throw result.error
                }

                const rows = result.data?.tenant?.transactions?.rows || []

                const rowCount = result.data?.tenant?.transactions.totalCount
                return {
                    rows: rows.map(r => ({ ...r, currency: r.currency.code })),
                    rowCount,
                }
            },
        }),
        [ fetchTransactions, tenantID, involvingAccountID, until ],
    )

    return (
        <div className="transactions-grid-content">
            <Snackbar open={fetchTransactionsResult.error !== undefined}
                      message={fetchTransactionsResult.error?.message} />
            <DataGridPremium<TransactionRow>
                disablePivoting
                disableRowGrouping
                disableAggregation
                onDataSourceError={e => console.error(e)}
                lazyLoading
                paginationModel={{ page: 0, pageSize: 30 }}
                dataSource={transactionsDataSource}
                columns={transactionColumns}
            />
        </div>
    )
}
