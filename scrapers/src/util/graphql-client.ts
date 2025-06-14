import { cacheExchange, Client, fetchExchange } from "@urql/core"

const graphQLURL = process.env.GRAPHQL_API_URL
if (!graphQLURL) {
    throw new Error("missing GRAPHQL_API_URL environment variable")
}

export const graphQLClient = new Client({
    url: graphQLURL,
    exchanges: [ cacheExchange, fetchExchange ],
})
