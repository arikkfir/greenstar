import { useMemo } from "react"

export function useTenantID(): string {
    return useMemo<string>((): string => {
        const match = /^([a-zA-Z0-9-_]+)\.app\.(greenstar.(?:test|kfirs\.com))$/.exec(window.location.hostname)
        if (!match) {
            throw new Error(`could not extract tenant ID from hostname: ${window.location.hostname}`)
        } else {
            return match[1]
        }
    }, [window.location.hostname])
}
