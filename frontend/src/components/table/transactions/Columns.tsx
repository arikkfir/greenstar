import {GridColDef} from "@mui/x-data-grid-premium";
import {Transaction} from "../../../client/transaction.ts";
import {useCallback, useContext, useMemo} from "react";
import {LocaleContext} from "../../../providers/LocaleProvider.tsx";
import {Account} from "../../../client/account.ts";

function useCurrencyFormatter() {
    const locale = useContext(LocaleContext)

    const currencyFormat = useMemo(() => {
        if (locale.language && locale.currency) {
            return new Intl.NumberFormat(locale.language, {style: 'currency', currency: locale.currency})
        } else {
            return new Intl.NumberFormat(navigator.language)
        }
    }, [locale, navigator.language])

    return useCallback((v: number) => (currencyFormat.format(v)), [currencyFormat])
}

function useDateFormatter() {
    const locale = useContext(LocaleContext)

    const dateTimeFormat = useMemo(() => {
        if (locale.language && locale.currency) {
            return new Intl.DateTimeFormat(locale.language, {dateStyle: 'short'})
        } else {
            return new Intl.DateTimeFormat(navigator.language, {dateStyle: 'short'})
        }
    }, [locale, navigator.language])

    return useCallback((v: number) => (dateTimeFormat.format(v)), [dateTimeFormat])
}

function useAccountFormatter(accountsByID: { [p: string]: Account }) {
    return useCallback((id: string) => {
        const account = accountsByID[id]
        return account ? account.displayName : id
    }, [accountsByID])
}

export function useColumns(accountsByID: { [p: string]: Account }) {
    const currencyFormatter = useCurrencyFormatter()
    const dateFormatter = useDateFormatter()
    const accountFormatter = useAccountFormatter(accountsByID)
    const columns: GridColDef<Transaction>[] = useMemo(() => ([
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: 'id',
                headerName: "ID",
                type: "string",
            },
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: "sourceAccountId", headerName: "From", type: "string", valueFormatter: accountFormatter,
            },
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: "date", headerName: "Date", type: "dateTime", valueFormatter: dateFormatter,
            },
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: "targetAccountId", headerName: "To", type: "string", valueFormatter: accountFormatter,
            },
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: "convertedAmount",
                headerName: "Amount",
                type: "number",
                valueFormatter: currencyFormatter,
            },
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: "description", headerName: "Description", type: "string",
            },
            {
                hideable: true, sortable: true, resizable: true, editable: false,
                groupable: false, pinnable: true, align: "left", filterable: true,
                field: "referenceId", headerName: "Reference", type: "string",
            },
        ]),
        [currencyFormatter, accountFormatter]
    );

    return columns
}
