import { gql } from "@urql/core"
import { graphQLClient } from "./graphql-client.js"
import { CreateScraperMutation, CreateScraperMutationVariables, Tenant } from "./graphql/graphql.js"
import { ACMEScraper } from "./main.js"
import { splitCamelCase } from "./util.js"

export const CreateScraper = gql(`
    mutation CreateScraper(
        $tenantID: ID!
        $scraperTypeID: ID!
        $id: ID!
        $displayName: String!
        $parameters: [ScraperParameterInput!]!
    ) {
        createScraper(
            tenantID: $tenantID
            scraperTypeID: $scraperTypeID
            id: $id
            displayName: $displayName
            parameters: $parameters
        ) {
            id
        }
    }
`)

export async function generateScraper(tenantID: Tenant["id"], scraper: ACMEScraper) {
    const result = await graphQLClient.mutation<CreateScraperMutation, CreateScraperMutationVariables>(
        CreateScraper,
        {
            tenantID,
            scraperTypeID: scraper.type,
            id: scraper.id,
            displayName: scraper.displayName || splitCamelCase(scraper.id),
            parameters: Object.entries(scraper.parameters || {}).map(kv => ({
                scraperTypeParameterID: kv[0],
                value: kv[1],
            })) || [],
        }).toPromise()
    if (result.error) {
        throw result.error
    }
}
