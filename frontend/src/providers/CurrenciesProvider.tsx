import { ReactNode, useContext, useEffect, useMemo } from "react"
import { LoadedCurrenciesStateCtx } from "../contexts/LoadingState.ts"
import { useQuery } from "@apollo/client"
import { CurrenciesContext } from "../contexts/Currencies.ts"
import { gql } from "../graphql"
import { delayFactor } from "./util.ts"

const CurrenciesQuery = gql(`
    query CurrenciesQuery {
        currencies {
            code
            createdAt
            updatedAt
            symbol
            nativeSymbol
            name
            namePlural
            decimalDigits
            countries
        }
    }
`)

export function CurrenciesProvider({ children }: { children: ReactNode }) {
    const { setLoaded: markAsLoaded }                             = useContext(LoadedCurrenciesStateCtx)
    const { data: currenciesResp, error: errorLoadingCurrencies } = useQuery(CurrenciesQuery)

    useEffect(() => {
        if (errorLoadingCurrencies) {
            markAsLoaded(errorLoadingCurrencies)
        } else if (currenciesResp?.currencies?.length) {
            window.setTimeout(markAsLoaded, Math.random() * delayFactor)
        } else if (currenciesResp && (!currenciesResp.currencies || !currenciesResp.currencies.length)) {
            markAsLoaded(new Error("No currencies loaded"))
        }
    }, [errorLoadingCurrencies, markAsLoaded, currenciesResp])

    const currencies = useMemo(
        () => (currenciesResp?.currencies || []),
        [currenciesResp?.currencies],
    )

    return (
        <CurrenciesContext.Provider value={currencies}>{children}</CurrenciesContext.Provider>
    )
}
