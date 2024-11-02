import {Transaction, useTransactionsClient} from "../../../client/transaction.ts";
import {Account} from "../../../client/account.ts";
import {DataGridPremiumProps, GridColDef} from "@mui/x-data-grid-premium";
import {useContext, useEffect, useMemo, useState} from "react";
import {LocaleContext} from "../../../providers/LocaleProvider.tsx";
import {useColumns} from "./Columns.tsx";
import {StatefulDataGrid} from "../StatefulDataGrid.tsx";

// const nativeDateTimeFormat = new Intl.DateTimeFormat(navigator.language, {
//     dateStyle: "short",
//     timeStyle: "short",
// });

// const nativeRelativeTimeFormat = new Intl.RelativeTimeFormat(navigator.language, {style: 'short'});

export interface DataGridProps extends Omit<DataGridPremiumProps<Transaction>, "columns" | "loading" | "rows"> {
}

export interface TransactionsTableProps {
    accounts: Account[]
    dataGridProps?: DataGridProps,
    sourceAccountId?: string,
    targetAccountId?: string,
}

export function TransactionsTable({accounts, dataGridProps, sourceAccountId, targetAccountId}: TransactionsTableProps) {
    const locale = useContext(LocaleContext)
    const transactionsClient = useTransactionsClient()
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<Error | undefined>()
    const [rowCount, setRowCount] = useState<number>(-1)
    const [transactions, setTransactions] = useState<Transaction[]>([])
    const [paginationModel, setPaginationModel] = useState({page: 0, pageSize: -1});

    const accountsByID = useMemo(() => {
        const accountsByID: { [key: string]: Account } = {}
        accounts.forEach(acc => accountsByID[acc.id] = acc)
        return accountsByID
    }, [accounts]);

    const columns: GridColDef<Transaction>[] = useColumns(accountsByID);

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
        <StatefulDataGrid<Transaction> stateId="transactions"
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
