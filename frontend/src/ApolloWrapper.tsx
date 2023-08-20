import React, {PropsWithChildren} from "react";
import {ApolloClient, ApolloProvider, InMemoryCache} from "@apollo/client";
import {gqlQueryURL} from "./config.ts";
import {useSession} from "@descope/react-sdk";

interface ApolloWrapperProps {
}

export function ApolloWrapper({children}: PropsWithChildren<ApolloWrapperProps>) {
    const {sessionToken} = useSession()

    const apolloClient = React.useMemo(() => new ApolloClient({
            cache: new InMemoryCache(),
            uri: gqlQueryURL,
            headers: {
                authorization: `Bearer ${sessionToken}`,
            },
            credentials: 'include',
        }),
        [sessionToken],
    )

    return (
        <ApolloProvider client={apolloClient}>
            {children}
        </ApolloProvider>
    )
}
