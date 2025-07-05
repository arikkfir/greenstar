import { gql } from "../../graphql"
import { FetchScrapersQuery } from "../../graphql/graphql.ts"

const FetchScraper = gql(`
    query FetchScraper($tenantID: ID!, $scraperID: ID!) {
        tenant(id: $tenantID) {
            scraper(id: $scraperID) {
                id
                createdAt
                updatedAt
                displayName
                lastSuccessfulScrapedDate
                type {
                    id
                    displayName
                    createdAt
                    updatedAt
                    parameters {
                        id
                        displayName
                        type
                    }
                }
                parameters {
                    parameter {
                        id
                        displayName
                        type
                    }
                    value
                }
            }
        }
    }
`)

export const FetchAccounts = gql(`
    query FetchAccounts($tenantID: ID!, $filter: String) {
        tenant(id: $tenantID) {
            accounts(filter: $filter) {
                id
                displayName
            }
        }
    }
`)

export const FetchScrapers = gql(`
    query FetchScrapers($tenantID: ID!) {
        tenant(id: $tenantID) {
            id
            scrapers {
                id
                createdAt
                updatedAt
                displayName
                lastSuccessfulScrapedDate
                type {
                    id
                    displayName
                    createdAt
                    updatedAt
                    parameters {
                        id
                        displayName
                        type
                    }
                }
                parameters {
                    parameter {
                        id
                        displayName
                        type
                    }
                    value
                }
            }
        }
    }
`)

export type ScraperRow = NonNullable<FetchScrapersQuery["tenant"]> [ "scrapers" ][number]
