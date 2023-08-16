import {graphql} from "../gql";

export const rootAccounts = graphql(/* GraphQL */ `
    query rootAccounts ($tenantID: ID!) {
        tenant(id: $tenantID) {
            id
            accounts {
                id
                displayName
                childCount
                icon
            }
        }
    }
`)

export const accountChildren = graphql(/* GraphQL */ `
    query accountChildren ($tenantID: ID!, $accountID: ID!) {
        tenant(id: $tenantID) {
            id
            account(id: $accountID) {
                id
                children {
                    id
                    displayName
                    childCount
                    icon
                }
            }
        }
    }
`)
