import { Context, ReactNode, useCallback, useContext, useMemo, useState } from "react"
import { LoadedState } from "../contexts/LoadingState.tsx"
import { delayFactor } from "./util.ts"

export function LoadingStateProvider({ ctx, children }: { ctx: Context<LoadedState>, children: ReactNode }) {
    const value: LoadedState  = useContext(ctx)
    const [loaded, setLoaded] = useState<boolean>(false)
    const [error, setError]   = useState<Error | undefined>()

    const _setLoaded = useCallback(
        (error?: Error): void => {
            if (error) {
                setError(error)
            } else if (!loaded) {
                window.setTimeout(() => setLoaded(true), Math.random() * delayFactor)
                setError(undefined)
            }
        },
        [loaded, setLoaded, setError],
    )

    const actual = useMemo(
        () => ({ name: value.name, loaded, error, setLoaded: _setLoaded }),
        [value.name, loaded, error, _setLoaded],
    )

    return (
        <ctx.Provider value={actual}>{children}</ctx.Provider>
    )
}
