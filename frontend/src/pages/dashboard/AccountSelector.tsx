import { Autocomplete, TextField } from "@mui/material"
import { useCallback, useContext, useEffect, useMemo, useState } from "react"
import { useQuery } from "@apollo/client"
import { gql } from "../../graphql"
import { useTenantID } from "../../hooks/tenant.ts"
import { SelectedCurrencyContext } from "../../contexts/SelectedCurrency.ts"

const AllAccounts = gql(`
    query AllAccountsQuery($tenantID: ID!, $currency: String!) {
        tenant(id: $tenantID) {
            id
            accounts {
                id
                displayName
                icon
                balance(currency: $currency)
            }
        }
    }
`)

export interface AccountData {
    id: string;
    displayName: string;
    icon: string;
    balance: number;
}

interface AccountSelectorProps {
    storageKey: string;
    label: string;
    onChange: (account: AccountData | null) => void;
}

export function AccountSelector({ storageKey, label, onChange }: AccountSelectorProps) {
    const tenantID                                = useTenantID()
    const selectedCurrency                        = useContext(SelectedCurrencyContext)
    const [ selectedAccount, setSelectedAccount ] = useState<AccountData | null>(null)

    // Fetch accounts - only once, no filter parameter
    const { data, loading, error } = useQuery(AllAccounts, {
        variables: {
            tenantID,
            currency: selectedCurrency.currency?.code || "USD",
        },
        skip: !selectedCurrency.currency,
    })

    // Memorized shortcut to refer to the accounts or an empty array if they haven't been fetched yet
    const accounts = useMemo(
        () => data?.tenant?.accounts || [],
        [ data?.tenant?.accounts ],
    )

    // Load the selected account from localStorage on initial render
    useEffect(
        () => {
            const storedAccountId = localStorage.getItem(storageKey)
            if (storedAccountId && accounts.length > 0) {
                const account = accounts.find(a => a.id === storedAccountId)
                if (account) {
                    setSelectedAccount(account)
                    onChange(account)
                }
            }
        },
        [ accounts, onChange, storageKey ],
    )

    // Handle account selection change
    const handleAccountChange = useCallback(
        (_event: any, newValue: AccountData | null) => {
            setSelectedAccount(newValue)
            onChange(newValue)
            if (newValue) {
                localStorage.setItem(storageKey, newValue.id)
            } else {
                localStorage.removeItem(storageKey)
            }
        },
        [ onChange, storageKey ],
    )


    if (loading) {
        return <TextField label={label} disabled value="Loading accounts..." />
    }
    if (error) {
        return <TextField label={label} disabled error value="Error loading accounts" />
    }

    return (
        <Autocomplete
            className="account-selector"
            options={accounts}
            getOptionLabel={(option) => option.displayName}
            value={selectedAccount}
            onChange={handleAccountChange}
            renderInput={(params) => <TextField {...params} label={label} />}
            isOptionEqualToValue={(option, value) => option.id === value.id}
        />
    )
}
