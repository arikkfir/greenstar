import { gql } from "@urql/core"
import { graphQLClient } from "./graphql-client.js"
import { Account, CreateAccountMutation, CreateAccountMutationVariables, Tenant } from "./graphql/graphql.js"
import { ACMEAccount } from "./main.js"
import { splitCamelCase } from "./util.js"

export const CreateAccount = gql(`
    mutation CreateAccount(
        $tenantID: ID!
        $id: ID!
        $parentID: ID
        $displayName: String!
        $icon: String
        $type: AccountType
    ) {
        createAccount(
            tenantID: $tenantID
            id: $id
            parentID: $parentID
            icon: $icon
            type: $type
            displayName: $displayName
        ) {
            id
            displayName
        }
    }
`)

export async function generateAccount(tenantID: Tenant["id"], builtinAccountIDs: Account["id"][], account: ACMEAccount,parent?: ACMEAccount) {
    if (!builtinAccountIDs.includes(account.id)) {
        const { id, displayName, icon } = account
        const creationResult = await graphQLClient.mutation<CreateAccountMutation, CreateAccountMutationVariables>(
            CreateAccount,
            {
                id,
                displayName: displayName || splitCamelCase(account.id),
                icon,
                type: null,
                parentID: parent?.id || null,
                tenantID
            }).toPromise()
        if (creationResult.error) {
            throw creationResult.error
        }
    }

    if (account.children?.length) {
        for (let child of account.children) {
            await generateAccount(tenantID, builtinAccountIDs, child, account)
        }
    }
}
