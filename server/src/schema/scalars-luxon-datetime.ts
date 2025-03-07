import type { GraphQLScalarTypeConfig } from "graphql"
import { GraphQLScalarType, Kind } from "graphql"
import { DateTime } from "luxon"

const config: GraphQLScalarTypeConfig<DateTime<true>, DateTime<true>> = {
    name: "DateTime",
    description:
        "A date-time string at UTC, such as 2007-12-03T10:15:30Z, " +
        "compliant with the `date-time` format outlined in section 5.6 of " +
        "the RFC 3339 profile of the ISO 8601 standard for representation " +
        "of dates and times using the Gregorian calendar.",
    serialize(value): DateTime<true> {
        if (value instanceof DateTime) {
            if (value.isValid) {
                return value
            } else {
                throw new Error(
                    `invalid given DateTime '${value}': ${value.invalidReason}, ${value.invalidExplanation}`,
                )
            }
        } else if (value instanceof Date) {
            const dt = DateTime.fromJSDate(value)
            if (dt.isValid) {
                return dt
            } else {
                throw new Error(`invalid given DateTime string: ${value}`)
            }
        } else if (typeof value === "string") {
            const dt = DateTime.fromISO(value, { setZone: true })
            if (dt.isValid) {
                return dt
            } else {
                throw new Error(`invalid given DateTime string: ${value}`)
            }
        } else {
            throw new Error(`unsupported value ${typeof value}: ${JSON.stringify(value, null, 2)}`)
        }
    },
    parseValue(value) {
        if (value instanceof DateTime) {
            if (value.isValid) {
                return value
            } else {
                throw new Error(`invalid date: ${value}`)
            }
        } else if (typeof value === "string") {
            const dt = DateTime.fromISO(value, { setZone: true })
            if (dt.isValid) {
                return dt
            } else {
                throw new Error(`invalid date: ${value}`)
            }
        } else {
            throw new Error(`unsupported value: ${JSON.stringify(value)}`)
        }
    },
    parseLiteral(ast) {
        if (ast.kind !== Kind.STRING) {
            throw new Error(`DateTime cannot represent non string or Date type ${"value" in ast && ast.value}`)
        }
        const { value } = ast
        const dt        = DateTime.fromISO(value, { setZone: true })
        if (dt.isValid) {
            return dt
        } else {
            throw new Error(`invalid date ${String(value)}`)
        }
    },
    extensions: {
        codegenScalarType: "DateTime | string",
        jsonSchema: {
            type: "string",
            format: "date-time",
        },
    },
}

export const DateTimeScalar = new GraphQLScalarType(config)
