import {useCallback, useContext, useEffect, useState} from "react";
import {Account, useAccountsClient} from "../client/account.ts";
import {Paper} from "@mui/material";
import Grid from "@mui/material/Grid2";
import {AccountsTree} from "../components/tree/accounts/AccountsTree.tsx";
import {TransactionsTable} from "../components/table/transactions/TransactionsTable.tsx";
import {LocaleContext} from "../providers/LocaleProvider.tsx";

export function TransactionsPage() {
    const locale = useContext(LocaleContext)
    const accountsClient = useAccountsClient()
    const [loadingAccounts, setLoadingAccounts] = useState(true)
    const [error, setError] = useState<Error | undefined>()
    const [accounts, setAccounts] = useState<Account[]>([])
    const [selectedAccount, setSelectedAccount] = useState<Account | undefined>()

    const handleEscape = useCallback((event: KeyboardEvent) => {
        if (event.key === 'Escape') {
            setSelectedAccount(undefined);
        }
    }, [setSelectedAccount]);

    useEffect(() => {
        window.addEventListener('keydown', handleEscape);
        return () => window.removeEventListener('keydown', handleEscape);
    }, []);

    useEffect(() => {
        if (locale.currency) {
            setLoadingAccounts(true)
            setError(undefined)
            accountsClient.List({currency: locale.currency})
                .then(r => {
                    setAccounts(r.items)
                    setSelectedAccount(undefined)
                })
                .catch(e => setError(e))
                .finally(() => setLoadingAccounts(false))
        } else if (locale.error) {
            setLoadingAccounts(false)
            setAccounts([])
            setError(locale.error)
        } else {
            setLoadingAccounts(false)
            setAccounts([])
            setError(undefined)
        }
    }, [locale, setLoadingAccounts, setError, accountsClient, setAccounts, setSelectedAccount]);

    if (error) {
        return (
            <p>Error : {error.message}</p>
        )
    }

    return (
        <Grid container spacing={2} sx={{width: '100%'}}>
            <Grid size={3}>
                <Paper sx={{height: '100%'}}>
                    <AccountsTree loading={loadingAccounts}
                                  accounts={accounts}
                                  selectedAccount={selectedAccount}
                                  onAccountSelected={acc => setSelectedAccount(acc)}/>
                </Paper>
            </Grid>
            <Grid size={9}>
                <Paper sx={{height: '100%', minWidth: 0, minHeight: 0}}>
                    <TransactionsTable accounts={accounts}
                                       sourceAccountId={selectedAccount?.id}
                                       targetAccountId={selectedAccount?.id}/>
                </Paper>
            </Grid>
        </Grid>
    )
}
