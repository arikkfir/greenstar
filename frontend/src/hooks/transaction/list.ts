// This file is generated by greenstar scripts. DO NOT EDIT.
// noinspection DuplicatedCode

import { Transaction } from "../../models/transaction.ts"
import { DateProperties } from "../../models/transaction.ts"
import { BaseAPIURL, QueryNilValue } from "../../services/util.ts"
import { useOperation, Method, Hook } from "../../util/operation.ts"

export type SortableProperty =
    | "amount"
    | "convertedAmount"
    | "currency"
    | "date"
    | "description"
    | "referenceId"
    | "sourceAccountId"
    | "targetAccountId"

export interface Request {
    offset?: number
    count?: number
    sort?: SortableProperty[]
    description?: string
    maxAmount?: number
    maxDate?: Date
    minAmount?: number
    minDate?: Date
    referenceId?: string
    sourceAccountId?: string
    targetAccountId?: string
    currency: string
}
export interface Response {
    offset: number
    totalCount: number
    items: Transaction[]
}

function buildURL(req: Request): string {
    const urlParams = new URLSearchParams()

    if (req.offset !== undefined && req.offset >= 0) {
        urlParams.set("_offset", req.offset + "")
    }
    if (req.count !== undefined && req.count >= 0) {
        urlParams.set("_count", req.count + "")
    }
    if (req.sort !== undefined && req?.sort.length) {
        req.sort.forEach((s) => urlParams.append("_sort", s))
    }
    if (req.description !== undefined) {
        if (req.description === null) {
            urlParams.set("description", QueryNilValue)
        } else {
            urlParams.set("description", req.description)
        }
    }
    if (req.maxAmount !== undefined) {
        if (req.maxAmount === null) {
            urlParams.set("maxAmount", QueryNilValue)
        } else {
            urlParams.set("maxAmount", req.maxAmount + "")
        }
    }
    if (req.maxDate !== undefined) {
        if (req.maxDate === null) {
            urlParams.set("maxDate", QueryNilValue)
        } else {
            urlParams.set("maxDate", req.maxDate.toISOString())
        }
    }
    if (req.minAmount !== undefined) {
        if (req.minAmount === null) {
            urlParams.set("minAmount", QueryNilValue)
        } else {
            urlParams.set("minAmount", req.minAmount + "")
        }
    }
    if (req.minDate !== undefined) {
        if (req.minDate === null) {
            urlParams.set("minDate", QueryNilValue)
        } else {
            urlParams.set("minDate", req.minDate.toISOString())
        }
    }
    if (req.referenceId !== undefined) {
        if (req.referenceId === null) {
            urlParams.set("referenceId", QueryNilValue)
        } else {
            urlParams.set("referenceId", req.referenceId)
        }
    }
    if (req.sourceAccountId !== undefined) {
        if (req.sourceAccountId === null) {
            urlParams.set("sourceAccountId", QueryNilValue)
        } else {
            urlParams.set("sourceAccountId", req.sourceAccountId)
        }
    }
    if (req.targetAccountId !== undefined) {
        if (req.targetAccountId === null) {
            urlParams.set("targetAccountId", QueryNilValue)
        } else {
            urlParams.set("targetAccountId", req.targetAccountId)
        }
    }
    if (req.currency !== undefined) {
        if (req.currency === null) {
            urlParams.set("currency", QueryNilValue)
        } else {
            urlParams.set("currency", req.currency)
        }
    }

    return `${BaseAPIURL}/transactions${urlParams.size ? "?" + urlParams.toString() : ""}`
}

export function useListTransactions(): Hook<Request, Response> {
    const opts = {
        initial: { offset: 0, totalCount: 0, items: [] },
        method: "GET" as Method,
        url: buildURL,
        dateProperties: DateProperties,
    }
    return useOperation<Request, Response>(opts)
}