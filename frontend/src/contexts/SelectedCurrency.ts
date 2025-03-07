import { createContext } from "react"
import { Currency } from "../graphql/graphql.ts"

export type SelectedCurrency = {
    currency?: Currency
    setCurrency: (c: Currency) => void
}

export const SelectedCurrencyContext = createContext<SelectedCurrency>({
    setCurrency: () => {
        throw new Error("setCurrency called on uninitialized SelectedCurrencyContext.")
    },
})
