import {useDomain} from "../hooks/domain.ts";
import {useTenantID as newUseTenantID} from "../hooks/tenant.ts";

export function useTenantID(): string {
    // TODO: migrate generated client code to use the new "useTenantID()"
    return newUseTenantID()
}

// TODO: move BaseAPIURL into a hook
export const BaseAPIURL = `https://api.${useDomain()}`

// TODO: move QueryNilValue into a hook
export const QueryNilValue = "<nil>"
