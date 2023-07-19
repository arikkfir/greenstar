import React, {PropsWithChildren} from "react";
import {ApolloClient, ApolloProvider, InMemoryCache} from "@apollo/client";
import {gqlQueryURL} from "./config";
import {useSession} from "@descope/react-sdk";

const TenantIDHeader = "x-greenstar-tenant-id"

interface ApolloWrapperProps {
    tenant: string
}

export function ApolloWrapper({tenant, children}: PropsWithChildren<ApolloWrapperProps>) {
    const {sessionToken} = useSession()

    const apolloClient = React.useMemo(() => new ApolloClient({
            cache: new InMemoryCache(),
            uri: gqlQueryURL,
            headers: {
                authorization: `Bearer ${sessionToken}`,
                [TenantIDHeader]: tenant,
            },
            credentials: 'include',
            defaultOptions: {
                watchQuery: {
                    fetchPolicy: 'network-only',
                    nextFetchPolicy: 'cache-only'
                },
            },
        }),
        [sessionToken, tenant],
    )

    return (
        <ApolloProvider client={apolloClient}>
            {children}
        </ApolloProvider>
    )
}
