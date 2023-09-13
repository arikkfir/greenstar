import {EnvironmentHelper} from "./env";
import {getSdk} from "../generated/graphql";
import {GraphQLClient} from "graphql-request";
import * as https from "https";

export interface Tenant {
    id: string
    displayName: string
}

export class BackendHelper {
    private readonly graphClient: GraphQLClient;

    constructor(private readonly env: EnvironmentHelper) {
        this.graphClient = new GraphQLClient(this.env.apiGraphURL, {
            agent: new https.Agent({rejectUnauthorized: false}),
            cache: "no-cache",
            method: "POST",
        })
    }

    async createTenant(token: string, id: string, name: string): Promise<Tenant> {
        const sdk = getSdk(this.graphClient)
        const response = await sdk.createTenant({
            tenantID: id,
            displayName: name,
        }, {"Authorization": `Bearer ${token}`})

        return response.createTenant
    }

    async loadTenant(token: string, id: string): Promise<Tenant> {
        const sdk = getSdk(this.graphClient)
        const response = await sdk.loadTenant({tenantID: id}, {"Authorization": `Bearer ${token}`})
        if (!response.tenant) {
            throw new Error(`Empty tenant received`)
        }
        return response.tenant
    }

    async deleteTenant(token: string, id: string) {
        const sdk = getSdk(this.graphClient)
        await sdk.deleteTenant({tenantID: id}, {"Authorization": `Bearer ${token}`})
    }
}
