import {Account} from "../../../client/account.ts";
import {Transaction, useTransactionsClient} from "../../../client/transaction.ts";
import {DataGridPremiumProps, GridColDef, GridColumnVisibilityModel} from "@mui/x-data-grid-premium";
import {LocaleContext} from "../../../providers/LocaleProvider.tsx";
import {useCurrencyFormatter, useDateFormatter} from "../../../hooks/locale.tsx";
import {useCallback, useContext, useEffect, useMemo, useState} from "react";
import {MapAccountsByID} from "../../../client/account-util.ts";
import {CustomDataGrid} from "../CustomDataGrid.tsx";
import type {GridAutosizeOptions} from "@mui/x-data-grid/hooks/features/columnResize";

export interface TransactionsTableProps extends Omit<DataGridPremiumProps<Transaction>, "columns" | "loading" | "rows"> {
    stateId: string
    accounts: Account[]
    sourceAccountId?: string,
    targetAccountId?: string,
}

export function TransactionsTable(props: TransactionsTableProps) {
    const {stateId, accounts, sourceAccountId, targetAccountId, ...dataGridProps} = props
    const locale = useContext(LocaleContext)
    const transactionsClient = useTransactionsClient()
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<Error | undefined>()
    const [rowCount, setRowCount] = useState<number>(-1)
    const [transactions, setTransactions] = useState<Transaction[]>([])
    const [paginationModel, setPaginationModel] = useState({page: 0, pageSize: -1});
    const [colVisibilityModel, setColVisibilityModel] = useState<GridColumnVisibilityModel>({
        "id": false,
        "referenceId": false,
    })
    const accountsByID = useMemo(() => MapAccountsByID(accounts), [accounts]);

    const currFmt = useCurrencyFormatter()
    const dateFmt = useDateFormatter()
    const accFmt = useCallback((id: string) => accountsByID[id]?.displayName || id, [accountsByID])
    const columns: GridColDef<Transaction>[] = useMemo(() => ([
        {field: "sourceAccountId", headerName: "From", type: "string", minWidth: 200, flex: 0, valueFormatter: accFmt},
        {field: "date", headerName: "Date", type: "dateTime", minWidth: 80, flex: 0, groupable: true, valueFormatter: dateFmt},
        {field: "targetAccountId", headerName: "To", type: "string", minWidth: 200, flex: 0, groupable: true, valueFormatter: accFmt},
        {field: "convertedAmount", headerName: "Amount", type: "number", minWidth: 120, flex: 0, valueFormatter: currFmt},
        {field: "description", headerName: "Description", type: "string", flex: 1},
        {field: "referenceId", headerName: "Reference", type: "string", flex: 0},
    ]), [accountsByID]);

    useEffect(() => {
        if (locale.currency && paginationModel.pageSize > 0) {
            setLoading(true)
            setError(undefined)
            const offset = paginationModel.page * paginationModel.pageSize
            const count = paginationModel.pageSize
            transactionsClient.List({offset, count, sourceAccountId, targetAccountId, currency: locale.currency})
                .then(r => {
                    setTransactions(r.items)
                    setRowCount(r.totalCount)
                })
                .catch(e => setError(e))
                .finally(() => setLoading(false))
        }
    }, [locale, paginationModel, setLoading, setError, transactionsClient, sourceAccountId, targetAccountId, setTransactions, setRowCount]);

    const autosizeOptions: GridAutosizeOptions = useMemo(() => ({
        columns: ['sourceAccountId', 'date', 'targetAccountId', 'convertedAmount', 'referenceId'],
        includeHeaders: true,
        includeOutliers: false,
        expand: true,
    }), [])

    if (error) {
        return (
            <p>Error : {error.message}</p>
        )
    }

    return (
        <CustomDataGrid<Transaction> stateId={stateId}
                                     columnVisibilityModel={colVisibilityModel}
                                     onColumnVisibilityModelChange={(model) => setColVisibilityModel(model)}
                                     columns={columns}
                                     autoPageSize={true}
                                     pagination={true}
                                     paginationMode="server"
                                     paginationModel={paginationModel}
                                     onPaginationModelChange={setPaginationModel}
                                     autosizeOnMount={true}
                                     autosizeOptions={autosizeOptions}
                                     loading={loading}
                                     rows={transactions}
                                     rowCount={rowCount}
                                     {...dataGridProps} />
    )
}
