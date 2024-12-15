import { Transaction } from "../../../models/transaction.ts"
import { DataGridPremiumProps, GridColDef, GridColumnVisibilityModel } from "@mui/x-data-grid-premium"
import { LocaleContext } from "../../../providers/LocaleProvider.tsx"
import { useCurrencyFormatter, useDateFormatter } from "../../../hooks/locale.ts"
import { useCallback, useContext, useEffect, useMemo, useState } from "react"
import { CustomDataGrid } from "../CustomDataGrid.tsx"
import type { GridAutosizeOptions } from "@mui/x-data-grid/hooks/features/columnResize"
import { AccountID, AccountsMap } from "../../../models/account-addons.ts"
import { useListTransactions } from "../../../hooks/transaction/list.ts"

const defaultColumnVisibilityModel = {
    id: false,
    referenceId: false,
}

function useColumns(accounts: AccountsMap): GridColDef<Transaction>[] {
    const currFmt = useCurrencyFormatter()
    const dateFmt = useDateFormatter()
    const accFmt = useCallback((id: string) => accounts[id]?.displayName || id, [accounts])
    return useMemo(
        () => [
            {
                field: "sourceAccountId",
                headerName: "From",
                type: "string",
                minWidth: 200,
                flex: 0,
                valueFormatter: accFmt,
            },
            {
                field: "date",
                headerName: "Date",
                type: "dateTime",
                minWidth: 80,
                flex: 0,
                groupable: true,
                valueFormatter: dateFmt,
            },
            {
                field: "targetAccountId",
                headerName: "To",
                type: "string",
                minWidth: 200,
                flex: 0,
                groupable: true,
                valueFormatter: accFmt,
            },
            {
                field: "convertedAmount",
                headerName: "Amount",
                type: "number",
                minWidth: 120,
                flex: 0,
                valueFormatter: currFmt,
            },
            { field: "description", headerName: "Description", type: "string", flex: 1 },
            { field: "referenceId", headerName: "Reference", type: "string", flex: 0 },
        ],
        [accounts, accFmt, dateFmt, currFmt],
    )
}

export interface TransactionsTableProps
    extends Omit<DataGridPremiumProps<Transaction>, "columns" | "loading" | "rows"> {
    stateId: string
    accounts: AccountsMap
    sourceAccountId?: AccountID
    targetAccountId?: AccountID
}

export function TransactionsTable({
    stateId,
    accounts,
    sourceAccountId,
    targetAccountId,
    ...dataGridProps
}: TransactionsTableProps) {
    const locale = useContext(LocaleContext)
    const [paginationModel, setPaginationModel] = useState({ page: 0, pageSize: -1 })
    const [colVisibilityModel, setColVisibilityModel] =
        useState<GridColumnVisibilityModel>(defaultColumnVisibilityModel)
    const columns: GridColDef<Transaction>[] = useColumns(accounts)
    const { request: fetchTransactions, operation: transactionsStatus } = useListTransactions()

    useEffect(() => {
        const offset = paginationModel.page * paginationModel.pageSize
        const count = paginationModel.pageSize
        if (count > 0) {
            fetchTransactions({ offset, count, sourceAccountId, targetAccountId, currency: locale.currency })
        }
    }, [paginationModel, sourceAccountId, targetAccountId, locale])

    const autosizeOptions: GridAutosizeOptions = useMemo(
        () => ({
            columns: ["sourceAccountId", "date", "targetAccountId", "convertedAmount", "referenceId"],
            includeHeaders: true,
            includeOutliers: false,
            expand: true,
        }),
        [],
    )

    if (transactionsStatus.error) {
        return <p>Error : {transactionsStatus.error.message}</p>
    }

    return (
        <CustomDataGrid<Transaction>
            stateId={stateId}
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
            loading={transactionsStatus.inflight}
            rows={transactionsStatus.response?.items}
            rowCount={transactionsStatus.response?.totalCount}
            {...dataGridProps}
        />
    )
}
