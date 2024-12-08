import {SyntheticEvent, useCallback, useContext, useEffect, useState} from "react";
import {Account, useAccountsClient} from "../client/account.ts";
import {LocaleContext} from "../providers/LocaleProvider.tsx";
import {Box, Paper} from "@mui/material";
import {AccountsTree} from "../components/tree/accounts/AccountsTree.tsx";
import {TransactionsTable} from "../components/table/transactions/TransactionsTable.tsx";
import {AccountID} from "../client/account-util.ts";

export function TransactionsPage() {
    const locale = useContext(LocaleContext)
    const accountsClient = useAccountsClient()
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
            setError(undefined)
            accountsClient.List({currency: locale.currency})
                .then(r => {
                    setAccounts(r.items)
                    setSelectedAccount(undefined)
                })
                .catch(e => setError(e))
        } else if (locale.error) {
            setAccounts([])
            setError(locale.error)
        } else {
            setAccounts([])
            setError(undefined)
        }
    }, [locale, setError, accountsClient, setAccounts, setSelectedAccount]);

    const handleAccountSelectionChange = useCallback(
        (_: SyntheticEvent, id: AccountID | null) => setSelectedAccount(accounts.find(a => a.id === id)),
        [accounts, setSelectedAccount]
    )

    if (error) {
        return (
            <p>Error : {error.message}</p>
        )
    }

    return (
        <Box sx={{
            flexGrow: 1, flexShrink: 1,
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'center',
            alignItems: 'stretch',
            alignContent: 'stretch',
            gap: '1rem',
            p: 2, overflow: "hidden",
        }}>
            <Paper sx={{flexGrow: 0, flexShrink: 0, minWidth: '20rem', p: 1, overflow: 'scroll'}}>
                <AccountsTree accounts={accounts}
                              selectedAccount={selectedAccount}
                              onSelectedItemsChange={handleAccountSelectionChange}/>
            </Paper>
            <Paper sx={{
                flexGrow: 1, flexShrink: 1,
                display: 'flex', flexDirection: 'column',
                justifyContent: 'center', alignItems: 'stretch', alignContent: 'stretch',
                overflow: "hidden",
            }}>
                <TransactionsTable sx={{flexGrow: 1, flexShrink: 1}}
                                   stateId="transactions"
                                   accounts={accounts}
                                   sourceAccountId={selectedAccount?.id}
                                   targetAccountId={selectedAccount?.id}/>
            </Paper>
        </Box>
    )
}
