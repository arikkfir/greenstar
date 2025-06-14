import { ApolloClient, HttpLink, InMemoryCache, InMemoryCacheConfig, NormalizedCacheObject } from "@apollo/client"
import { luxonDateTimeMiddleware } from "./LuxonApolloLinkMiddleware.ts"

const graphQLURL = import.meta.env.VITE_GRAPHQL_API_URL
if (!graphQLURL) {
    throw new Error("missing VITE_GRAPHQL_API_URL environment variable")
}

const apolloCacheConfig: InMemoryCacheConfig = {
    typePolicies: {
        Tenant: {
            fields: {
                transactions: {
                    keyArgs: [ "tenantID", "involvingAccountID", "sort" ],
                    merge(existing, incoming) {
                        const existingRows = existing?.rows || []
                        const incomingRows = incoming?.rows || []
                        const offset       = incoming?.offset || existingRows.length
                        return existingRows.slice(0, offset).concat(incomingRows).concat(incomingRows)
                    },
                },
            },
        },
    },
}

export const apolloClient = new ApolloClient<NormalizedCacheObject>({
    link: luxonDateTimeMiddleware.concat(new HttpLink({ uri: graphQLURL })),
    cache: new InMemoryCache(apolloCacheConfig),
})
