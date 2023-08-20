import {PropsWithChildren} from "react";
import {getJwtPermissions, getSessionToken} from "@descope/react-sdk";

interface RequiresTenantPermissionProps {
    tenant: string
    anyOf: string[]
}

export function RequiresTenantPermission({tenant, anyOf, children}: PropsWithChildren<RequiresTenantPermissionProps>) {
    const permissions = getJwtPermissions(getSessionToken(), tenant);
    if (anyOf.some(p => permissions.includes(p))) {
        return children
    } else {
        return <></>
    }
}
