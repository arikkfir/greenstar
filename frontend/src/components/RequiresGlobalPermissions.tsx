import {PropsWithChildren} from "react";
import {getJwtPermissions, getSessionToken} from "@descope/react-sdk";
import {globalTenant} from "../config.ts";

interface RequiresGlobalPermissionsProps {
    anyOf: string[]
}

export function RequiresGlobalPermissions({anyOf, children}: PropsWithChildren<RequiresGlobalPermissionsProps>) {
    const permissions = getJwtPermissions(getSessionToken(), globalTenant);
    if (anyOf.some(p => permissions.includes(p))) {
        return children
    } else {
        return <></>
    }
}
