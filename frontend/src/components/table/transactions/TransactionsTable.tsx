import {Account} from "../../../client/account.ts";
import {Transaction, useTransactionsClient} from "../../../client/transaction.ts";
import {
    DataGridPremiumProps,
    GridCallbackDetails,
    GridColDef,
    GridColumnVisibilityModel
} from "@mui/x-data-grid-premium";
import {LocaleContext} from "../../../providers/LocaleProvider.tsx";
import {useCurrencyFormatter, useDateFormatter} from "../../../hooks/locale.tsx";
import {useCallback, useContext, useEffect, useMemo, useState} from "react";
import {AccountsMap, MapAccountsByID} from "../../../client/account-util.ts";
import {CustomDataGrid} from "../CustomDataGrid.tsx";

function buildGridColumnDefinitions(accountsByID: AccountsMap): GridColDef<Transaction>[] {
    const currencyFormatter = useCurrencyFormatter()
    const dateFormatter = useDateFormatter()
    const accountFormatter = useCallback((id: string) => accountsByID[id]?.displayName || id, [accountsByID])
    return useMemo(() => ([
            {field: 'id', headerName: "ID", groupable: false, type: "string"},
            {field: "sourceAccountId", headerName: "From", type: "string", valueFormatter: accountFormatter},
            {field: "date", headerName: "Date", type: "dateTime", valueFormatter: dateFormatter},
            {field: "targetAccountId", headerName: "To", type: "string", valueFormatter: accountFormatter},
            {
                field: "convertedAmount",
                headerName: "Amount",
                groupable: false,
                type: "number",
                valueFormatter: currencyFormatter
            },
            {field: "description", headerName: "Description", groupable: false, type: "string"},
            {field: "referenceId", headerName: "Reference", groupable: false, type: "string"},
        ]),
        [currencyFormatter, accountFormatter]
    );
}

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
    const columns: GridColDef<Transaction>[] = buildGridColumnDefinitions(accountsByID);

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

    const autosizeOptions = useMemo(() => ({expand: true}), [])

    if (error) {
        return (
            <p>Error : {error.message}</p>
        )
    }

    return (
        <CustomDataGrid<Transaction> stateId={stateId}
                                     columnVisibilityModel={colVisibilityModel}
                                     onColumnVisibilityModelChange={(model: GridColumnVisibilityModel, _: GridCallbackDetails) => setColVisibilityModel(model)}
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
