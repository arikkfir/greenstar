import {EnvironmentHelper} from "./env";
import {getSdk} from "../generated/graphql";
import {GraphQLClient} from "graphql-request";
import * as https from "https";

export class BackendHelper {
    private readonly _httpsAgent: https.Agent = new https.Agent({rejectUnauthorized: false})

    constructor(private readonly env: EnvironmentHelper) {
    }

    protected graphClient(token: string): GraphQLClient {
        return new GraphQLClient(this.env.apiGraphURL, {
            agent: this._httpsAgent,
            cache: "no-cache",
            headers: {"Authorization": `Bearer ${token}`},
            method: "POST",
        })
    }

    protected sdk(token: string) {
        return getSdk(this.graphClient(token))
    }
}
