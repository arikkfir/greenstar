import { gql } from "../../graphql"

export const FetchScraperTypes = gql(`
    query ScraperTypes {
        scraperTypes {
            id
            displayName
            parameters {
                id
                displayName
                type
            }
        }
    }
`)
