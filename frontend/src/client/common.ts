let tenantID = "";

export function useDomain(): string {
    const match = /^(?:[a-zA-Z0-9-_]+\.)?app\.(greenstar.(?:local|kfirs\.com))$/.exec(window.location.hostname)
    if (!match) {
        throw new Error(`could not extract domain from hostname: ${window.location.hostname}`)
    } else {
        return match[1]
    }
}

export function isHealthEndpoint(): boolean {
    return window.location.pathname === "/healthz"
}

export function isSignupPage(): boolean {
    return window.location.hostname == "app." + useDomain() && window.location.pathname.startsWith("/signup")
}

export function isTenantSelectionPage(): boolean {
    return window.location.hostname == "app." + useDomain() && window.location.pathname.startsWith("/tenants")
}

export function useTenantID(): string {
    if (tenantID == "") {
        const match = /^([a-zA-Z0-9-_]+)\.app\.(greenstar.(?:local|kfirs\.com))$/.exec(window.location.hostname)
        if (!match) {
            throw new Error(`could not extract tenant ID from hostname: ${window.location.hostname}`)
        } else {
            tenantID = match[1]
        }
    }
    return tenantID;
}

export const BaseAPIURL = `https://api.${useDomain()}`
export const QueryNilValue = "<nil>"
