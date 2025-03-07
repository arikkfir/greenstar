import { Context, createContext } from "react"

export type LoadedState = {
    readonly name: string
    readonly loaded: boolean,
    readonly error?: Error,
    setLoaded(error?: Error): void,
}

function createInitialLoadedContext(name: string): Context<LoadedState> {
    return createContext<LoadedState>({
        name,
        loaded: false,
        setLoaded: (_error?: Error) => {
            throw new Error("setLoaded called on uninitialized loaded-state context.")
        },
    })
}

export const LoadedLanguagesStateCtx: Context<LoadedState>        = createInitialLoadedContext("languages")
export const LoadedSelectedLanguageStateCtx: Context<LoadedState>     = createInitialLoadedContext("languages")
export const LoadedCurrenciesStateCtx: Context<LoadedState>       = createInitialLoadedContext("currencies")
export const LoadedSelectedCurrencyStateCtx: Context<LoadedState> = createInitialLoadedContext("selectedCurrency")
