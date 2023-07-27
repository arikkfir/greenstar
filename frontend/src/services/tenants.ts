import {graphql} from "../gql";

const tenantIDHashLetters = "abcdefghijklmnopqrstuvwxyz"

export function randomTenantID(length: number): string {
    let hash = ""
    for (let i = 0; i < length; i++) {
        hash += tenantIDHashLetters.charAt(Math.floor(Math.random() * tenantIDHashLetters.length))
    }
    return hash
}

export const allTenantsQuery = graphql(/* GraphQL */ `
    query allTenants {
        tenants {
            id, displayName
        }
    }
`)

export const createTenantMutation = graphql(/* GraphQL */`
    mutation CreateTenant($id: String, $displayName: String!) {
        createTenant(tenantID: $id, tenant: {displayName: $displayName}) {
            id
            displayName
        }
    }
`)

export const deleteTenantMutation = graphql(/* GraphQL */`
    mutation DeleteTenant($id: String!) {
        deleteTenant(tenantID: $id)
    }
`)

export const updateTenantMutation = graphql(/* GraphQL */`
    mutation UpdateTenant($id: String!, $displayName: String!) {
        updateTenant(tenantID: $id, tenant: {displayName: $displayName}) {
            id
            displayName
        }
    }
`)
