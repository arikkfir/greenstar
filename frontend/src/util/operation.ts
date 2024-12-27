import { useTenantID } from "../hooks/tenant.ts"
import { useCallback, useMemo, useState } from "react"
import { BadRequestError, InternalError, NotAuthenticatedError } from "./errors.ts"

export interface Operation<Response> {
    inflight: boolean
    response: Response
    error?: Error
}

export interface Hook<Request, Response> {
    request: (req: Request) => void
    operation: Operation<Response>
}

export type Method = "POST" | "GET" | "PATCH" | "PUT" | "DELETE"
export type OperationURL<Request> = string | ((req: Request) => string)

export interface OperationProps<Request, Response> {
    initial: Response
    method: Method
    url: OperationURL<Request>
    dateProperties?: string[]
}

export function useOperation<Request, Response>({
    initial,
    method,
    url,
    dateProperties: dateProps = [],
}: OperationProps<Request, Response>): Hook<Request, Response> {
    const tenantID = useTenantID()
    const jsonReviver = useCallback(
        (key: string, value: any) => (dateProps.includes(key) ? new Date(value) : value),
        dateProps,
    )
    const [inflight, setInflight] = useState<boolean>(false)
    const [error, setError] = useState<Error | undefined>()
    const [response, setResponse] = useState<Response>(initial)
    const request = function (req: Request): void {
        if (inflight) {
            return
        }
        setInflight(true)
        setResponse(initial)
        setError(undefined)

        const headers = {
            "Content-Type": "application/json",
            "X-GreenSTAR-Tenant-ID": tenantID,
            "X-Request-ID": crypto.randomUUID(),
        }

        let resolvedURL: string
        if (typeof url == "function") {
            resolvedURL = url(req)
        } else {
            resolvedURL = url
        }

        const body = method == "POST" || method == "PUT" || method == "PATCH" ? JSON.stringify(req) : undefined

        fetch(resolvedURL, {
            method,
            headers,
            body,
        })
            .then((response) => {
                if (response.status >= 500) {
                    throw InternalError
                } else if (response.status == 401) {
                    throw NotAuthenticatedError
                } else if (response.status >= 400) {
                    response
                        .text()
                        .then((text) => {
                            setInflight(false)
                            setResponse(initial)
                            setError(new BadRequestError(response.status, text))
                        })
                        .catch((e) => {
                            console.error("Failed reading response", e)
                            setInflight(false)
                            setResponse(initial)
                            setError(InternalError)
                        })
                } else {
                    response
                        .text()
                        .then((text) => (text ? JSON.parse(text, jsonReviver) : {}))
                        .then((json) => {
                            setInflight(false)
                            setResponse(json)
                            setError(undefined)
                        })
                        .catch((e) => {
                            console.error("Failed reading or parsing response", e)
                            setInflight(false)
                            setResponse(initial)
                            setError(InternalError)
                        })
                }
            })
            .catch((e) => {
                if (e === NotAuthenticatedError) {
                    setInflight(false)
                    setResponse(initial)
                    setError(e)
                } else {
                    console.error(`Failed sending request to '${resolvedURL}`, e)
                    setInflight(false)
                    setResponse(initial)
                    setError(InternalError)
                }
            })
    }
    return useMemo(
        () => ({
            request,
            operation: {
                inflight,
                response,
                error,
            },
        }),
        [inflight, response, error],
    )
}
