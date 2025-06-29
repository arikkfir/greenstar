/**
 * Bank Yahav Configuration
 *
 * This module provides configuration management for the Bank Yahav scraper.
 * It loads required environment variables for authentication and identification.
 */

import { DateTime } from "luxon"

/**
 * Configuration class for Bank Yahav scraper
 *
 * Loads and provides access to required configuration values from environment variables.
 * All configuration values are read-only and initialized during construction.
 */
export class Config {
    /**
     * Tenant identifier in the system
     */
    public readonly tenantID: string

    /**
     * Bank account identifier in the system
     */
    public readonly accountID: string

    /**
     * Username for Bank Yahav authentication
     */
    public readonly username: string

    /**
     * Password for Bank Yahav authentication
     */
    public readonly password: string

    /**
     * PIN number for Bank Yahav authentication
     */
    public readonly pinno: string

    /**
     * Indicates whether to download transactions XLS (Excel Spreadsheet) files or not.
     */
    public readonly downloadXLS: boolean

    /**
     * The year of the time period to scrape.
     */
    public readonly startDate: DateTime

    /**
     * The month of the time period to scrape.
     */
    public readonly endDate: DateTime

    /**
     * Creates a new configuration instance
     *
     * Loads all required environment variables for Bank Yahav scraper.
     * Throws an error if any required environment variable is missing.
     */
    constructor() {
        this.tenantID  = requireStringEnvVar("TENANT_ID")
        this.accountID = requireStringEnvVar("ACCOUNT_ID")
        console.info(`Scraping account ${this.accountID} in tenant ${this.tenantID}`)

        this.username = requireStringEnvVar("BANK_YAHAV_USERNAME")
        this.password = requireStringEnvVar("BANK_YAHAV_PASSWORD")
        this.pinno    = requireStringEnvVar("BANK_YAHAV_PINNO")

        this.downloadXLS = booleanEnvVar("DOWNLOAD_XLS", false) || false
        if (this.downloadXLS) {
            console.info(`Downloading XLS files is enabled`)
        } else {
            console.info(`Downloading XLS files is disabled`)
        }

        this.startDate = requireDateTimeEnvVar("START_DATE")
        this.endDate   = requireDateTimeEnvVar("END_DATE")
        console.info(`Scraping period set to: ${this.startDate.toFormat("yyyy-MM-dd")}-${this.endDate.toFormat(
            "yyyy-MM-dd")}`)
    }
}

function requireStringEnvVar(name: string) {
    const value = stringEnvVar(name)
    if (!value) {
        throw new Error(`Missing environment variable: ${name}`)
    }
    return value
}

function stringEnvVar(name: string, defaultValue?: string) {
    return process.env[name] || defaultValue
}

function requireDateTimeEnvVar(name: string) {
    const value = dateTimeEnvVar(name)
    if (!value) {
        throw new Error(`Missing environment variable: ${name}`)
    }
    return value
}

function dateTimeEnvVar(name: string) {
    const value = process.env[name]
    if (!value) {
        return undefined
    } else {
        return DateTime.fromISO(value)
    }
}

function booleanEnvVar(name: string, defaultValue?: boolean) {
    const value = process.env[name]
    if (!value) {
        return defaultValue
    } else {
        return [ "true", "t", "yes", "y", "ok" ].includes(value.toLowerCase())
    }
}
