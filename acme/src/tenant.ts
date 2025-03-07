import { gql } from "@urql/core"
import { graphQLClient } from "./graphql-client.js"
import {
    Account,
    CreateTenantMutation,
    CreateTenantMutationVariables,
    DeleteTenantMutation,
    DeleteTenantMutationVariables,
    FindTenantAccountsQuery,
    FindTenantAccountsQueryVariables,
    FindTenantQuery,
    FindTenantQueryVariables,
} from "./graphql/graphql.js"

export const FindTenant = gql(`
    query FindTenant($id: ID!) {
        tenant(id: $id) {
            id
        }
    }
`)

export const DeleteTenant = gql(`
    mutation DeleteTenant($id: ID!) {
        deleteTenant(id: $id)
    }
`)

export const CreateTenant = gql(`
    mutation CreateTenant($id: ID!, $displayName: String!) {
        createTenant(id: $id, displayName: $displayName) {
            id
            displayName
        }
    }
`)

export const FindTenantAccounts = gql(`
    query FindTenantAccounts($id: ID!) {
        tenant(id: $id) {
            accounts {
                id
            }
        }
    }
`)

export async function generateTenant(id: string, displayName: string): Promise<Account["id"][]> {
    console.info(`Checking if tenant ${id} exists...`)
    const findTenantResult =
              await graphQLClient.query<FindTenantQuery, FindTenantQueryVariables>(FindTenant, { id }).toPromise()
    if (findTenantResult.error) {
        throw findTenantResult.error
    }

    if (findTenantResult.data?.tenant) {
        console.info(`Tenant ${id} exists - deleting...`)
        const deletionResult =
                  await graphQLClient.mutation<DeleteTenantMutation, DeleteTenantMutationVariables>(DeleteTenant, { id })
                                     .toPromise()
        if (deletionResult.error) {
            throw deletionResult.error
        }
        console.info(`Deleted tenant ${id}`)
    }

    console.info(`Sending tenant ${id} creation request...`)
    const creationResult =
              await graphQLClient.mutation<CreateTenantMutation, CreateTenantMutationVariables>(
                  CreateTenant, { id, displayName }
              ).toPromise()
    if (creationResult.error) {
        throw creationResult.error
    }

    console.info(`Fetching tenant ${id} built-in accounts...`)
    const findTenantAccountsResult =
              await graphQLClient.query<FindTenantAccountsQuery, FindTenantAccountsQueryVariables>(
                  FindTenantAccounts, { id }
              ).toPromise()
    if (findTenantAccountsResult.error) {
        throw findTenantAccountsResult.error
    }

    const builtinAccounts = findTenantAccountsResult.data?.tenant?.accounts
    return builtinAccounts?.map(a => a.id) || []
}
