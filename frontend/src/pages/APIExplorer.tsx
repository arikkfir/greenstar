import {Stack} from "@mui/material";
import {gqlQueryURL} from "../config";
import {createGraphiQLFetcher} from '@graphiql/toolkit';
import {GraphiQL} from "graphiql";
import {useSession} from "@descope/react-sdk";
import "graphiql/graphiql.min.css";

export function APIExplorer() {
    const {sessionToken} = useSession();
    const fetcher = createGraphiQLFetcher({
        url: gqlQueryURL,
        headers: {
            Accept: 'application/json',
            Authorization: `Bearer ${sessionToken}`,
        },
    });
    return (
        <Stack direction="row" sx={{flexGrow: 1, height: "100%"}}>
            <Stack direction="column" sx={{flexGrow: 1}}>
                <GraphiQL fetcher={fetcher}/>
            </Stack>
        </Stack>
    )
}
