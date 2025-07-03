import pg, { Pool } from "pg"
import { DateTime, Duration } from "luxon"
import postgresInterval from "postgres-interval"

const { types } = pg

function convertDateValues(value: string | null): (DateTime | null) {
    if (value == null) {
        return null
    } else {
        return DateTime.fromSQL(value, { zone: "UTC" })
    }
}

function convertIntervalValues(value: any | null): (Duration | null) {
    if (value == null) {
        return null
    }

    const parsed = postgresInterval(value)

    // Convert to Luxon Duration
    return Duration.fromObject({
        years: parsed.years || 0,
        months: parsed.months || 0,
        days: parsed.days || 0,
        hours: parsed.hours || 0,
        minutes: parsed.minutes || 0,
        seconds: parsed.seconds || 0,
        milliseconds: parsed.milliseconds || 0,
    })
}

function convertIntegerValues(value: string | null): (number | null) {
    if (value == null) {
        return null
    } else {
        return parseInt(value)
    }
}

function convertFloatValues(value: string | null): (number | null) {
    if (value == null) {
        return null
    } else {
        return parseFloat(value)
    }
}

types.setTypeParser(1082, convertDateValues)   // DATE
types.setTypeParser(1083, convertDateValues)   // TIME
types.setTypeParser(1266, convertDateValues)   // TIMETZ
types.setTypeParser(1114, convertDateValues)   // TIMESTAMP
types.setTypeParser(1184, convertDateValues)   // TIMESTAMPTZ
types.setTypeParser(1186, convertIntervalValues) // INTERVAL
types.setTypeParser(21, convertIntegerValues)   // INT2
types.setTypeParser(23, convertIntegerValues)   // INT4
types.setTypeParser(20, convertIntegerValues)   // INT8
types.setTypeParser(700, convertFloatValues)    // FLOAT4
types.setTypeParser(701, convertFloatValues)    // FLOAT8
types.setTypeParser(1700, convertFloatValues)   // NUMERIC

export const pgPool = new Pool({
    max: 5,
    min: 0,
    idleTimeoutMillis: 1000 * 60 * 5,
    statement_timeout: 1000 * 60,
    connectionTimeoutMillis: 1000 * 5,
    keepAliveInitialDelayMillis: 1000 * 5,
    query_timeout: 1000 * 120,
})

const convertToCamelCase = (str: string): string => {
    return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase())
              .replace(/Id$/, "ID")
}

export const convertObjectKeysToCamelCase = <T extends object>(obj: T): { [p: string]: any } => {
    const newObj: { [key: string]: any } = {}
    Object.entries(obj).forEach(([ key, value ]) => newObj[convertToCamelCase(key)] = value)
    return newObj
}
