export function useDomain(): string {
    const match = /^(?:[a-zA-Z0-9-_]+\.)?app\.(greenstar.(?:test|kfirs\.com))$/.exec(window.location.hostname)
    if (!match) {
        throw new Error(`could not extract domain from hostname: ${window.location.hostname}`)
    } else {
        return match[1]
    }
}
