import { gql } from "../../graphql"

export const UpsertScraper = gql(`
    mutation UpsertScraper(
        $tenantID: ID!,
        $scraperID: ID,
        $scraperTypeID: ID!,
        $displayName: String!,
        $parameters: [ScraperParameterInput!]!
    ) {
        upsertScraper(
            tenantID: $tenantID,
            id: $scraperID,
            scraperTypeID: $scraperTypeID,
            displayName: $displayName,
            parameters: $parameters
        ) {
            id
            createdAt
            updatedAt
            displayName
            type {
                id
                displayName
            }
            parameters {
                parameter {
                    id
                    displayName
                }
                value
            }
        }
    }
`)
