import { ReactNode, useCallback, useContext, useEffect, useMemo, useState } from "react"
import { LoadedSelectedCurrencyStateCtx } from "../contexts/LoadingState.ts"
import { Currency } from "../graphql/graphql.ts"
import { CurrenciesContext } from "../contexts/Currencies.ts"
import { SelectedCurrencyContext } from "../contexts/SelectedCurrency.ts"
import { delayFactor } from "./util.ts"

export function SelectedCurrencyProvider({ children }: { children: ReactNode }) {
    const { setLoaded: markAsLoaded } = useContext(LoadedSelectedCurrencyStateCtx)
    const storageKey                  = "currency"
    const currencies: Currency[]      = useContext(CurrenciesContext)
    const [ currency, setCurrency ]   = useState<Currency | undefined>()

    useEffect(() => {
        if (!currency) {
            const storedCode = localStorage.getItem(storageKey)
            if (storedCode) {
                const currency = currencies.find((c) => c.code === storedCode)
                if (currency) {
                    setCurrency(currency)
                }
            } else {
                setCurrency(currencies.find((c) => c.code === "USD"))
            }
            window.setTimeout(markAsLoaded, Math.random() * delayFactor)
        }
    }, [ currencies, setCurrency, markAsLoaded ])

    const saveCurrency = useCallback(
        (c: Currency) => {
            setCurrency(c)
            localStorage.setItem(storageKey, c.code)
        },
        [ setCurrency ],
    )

    const value = useMemo(
        () => ({ currency, setCurrency: saveCurrency }),
        [ currency, saveCurrency ],
    )

    return (
        <SelectedCurrencyContext.Provider value={value}>{children}</SelectedCurrencyContext.Provider>
    )
}
