// This file is generated by greenstar scripts. DO NOT EDIT.
// noinspection DuplicatedCode

import { Tenant } from "../../models/tenant.ts"
import { DateProperties } from "../../models/tenant.ts"
import { BaseAPIURL } from "../../services/util.ts"
import { useOperation, Method, Hook } from "../../util/operation.ts"

export interface Request {
    id: string
}
export type Response = Tenant | undefined

function buildURL(req: Request): string {
    const urlParams = new URLSearchParams()

    return `${BaseAPIURL}/tenants/${req.id}${urlParams.size ? "?" + urlParams.toString() : ""}`
}

export function useGetTenant(): Hook<Request, Response> {
    const opts = {
        initial: undefined,
        method: "GET" as Method,
        url: buildURL,
        dateProperties: DateProperties,
    }
    return useOperation<Request, Response>(opts)
}