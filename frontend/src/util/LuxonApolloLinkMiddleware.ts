import { ApolloLink } from "@apollo/client"
import { DateTime } from "luxon"

const RFC3399 = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}(?::\d{2})?)?$/

function transformStringDatesToLuxonDateTime(value: any): any {
    if (!value || typeof value !== "object") {
        return value
    }

    if (Array.isArray(value)) {
        return value.map(item => transformStringDatesToLuxonDateTime(item))
    }

    const transformed = { ...value }
    for (const [ key, val ] of Object.entries(transformed)) {
        if (typeof val === "string" && RFC3399.test(val)) {
            const dateTime: DateTime = DateTime.fromISO(val, { setZone: true, locale: navigator.language })
            if (dateTime.isValid) {
                transformed[key] = dateTime
            }
        } else if (val && typeof val === "object") {
            transformed[key] = transformStringDatesToLuxonDateTime(val)
        }
    }
    return transformed
}

export const luxonDateTimeMiddleware = new ApolloLink((operation, forward) => {
    return forward(operation).map(transformStringDatesToLuxonDateTime)
})
