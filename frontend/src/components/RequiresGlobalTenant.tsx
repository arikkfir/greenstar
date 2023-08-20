import {PropsWithChildren} from "react";
import {globalTenant} from "../config.ts";

interface RequiresGlobalTenantProps {
    tenant: string
}

export function RequiresGlobalTenant({tenant, children}: PropsWithChildren<RequiresGlobalTenantProps>) {
    if (tenant == globalTenant) {
        return children
    } else {
        return <></>
    }
}
