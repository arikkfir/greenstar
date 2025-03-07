import { createContext } from "react"
import { Currency } from "../graphql/graphql.ts"

export const CurrenciesContext = createContext<Currency[]>([])
