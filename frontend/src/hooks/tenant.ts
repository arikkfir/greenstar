let tenantID = ""

export function useTenantID(): string {
    if (tenantID == "") {
        const match = /^([a-zA-Z0-9-_]+)\.app\.(greenstar.(?:test|kfirs\.com))$/.exec(window.location.hostname)
        if (!match) {
            throw new Error(`could not extract tenant ID from hostname: ${window.location.hostname}`)
        } else {
            tenantID = match[1]
        }
    }
    return tenantID
}
