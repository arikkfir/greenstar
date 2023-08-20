/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n    query rootAccounts ($tenantID: ID!) {\n        tenant(id: $tenantID) {\n            id\n            accounts {\n                id\n                displayName\n                childCount\n                icon\n            }\n        }\n    }\n": types.RootAccountsDocument,
    "\n    query accountChildren ($tenantID: ID!, $accountID: ID!) {\n        tenant(id: $tenantID) {\n            id\n            account(id: $accountID) {\n                id\n                children {\n                    id\n                    displayName\n                    childCount\n                    icon\n                }\n            }\n        }\n    }\n": types.AccountChildrenDocument,
    "\n    query allTenants {\n        tenants {\n            id\n            displayName\n        }\n    }\n": types.AllTenantsDocument,
    "\n    mutation CreateTenant($id: String, $displayName: String!) {\n        tenant: createTenant(tenantID: $id, tenant: {displayName: $displayName}) {\n            id\n            displayName\n        }\n    }\n": types.CreateTenantDocument,
    "\n    mutation DeleteTenant($id: String!) {\n        id: deleteTenant(tenantID: $id)\n    }\n": types.DeleteTenantDocument,
    "\n    mutation UpdateTenant($id: String!, $displayName: String!) {\n        tenant: updateTenant(tenantID: $id, tenant: {displayName: $displayName}) {\n            id\n            displayName\n        }\n    }\n": types.UpdateTenantDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    query rootAccounts ($tenantID: ID!) {\n        tenant(id: $tenantID) {\n            id\n            accounts {\n                id\n                displayName\n                childCount\n                icon\n            }\n        }\n    }\n"): (typeof documents)["\n    query rootAccounts ($tenantID: ID!) {\n        tenant(id: $tenantID) {\n            id\n            accounts {\n                id\n                displayName\n                childCount\n                icon\n            }\n        }\n    }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    query accountChildren ($tenantID: ID!, $accountID: ID!) {\n        tenant(id: $tenantID) {\n            id\n            account(id: $accountID) {\n                id\n                children {\n                    id\n                    displayName\n                    childCount\n                    icon\n                }\n            }\n        }\n    }\n"): (typeof documents)["\n    query accountChildren ($tenantID: ID!, $accountID: ID!) {\n        tenant(id: $tenantID) {\n            id\n            account(id: $accountID) {\n                id\n                children {\n                    id\n                    displayName\n                    childCount\n                    icon\n                }\n            }\n        }\n    }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    query allTenants {\n        tenants {\n            id\n            displayName\n        }\n    }\n"): (typeof documents)["\n    query allTenants {\n        tenants {\n            id\n            displayName\n        }\n    }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    mutation CreateTenant($id: String, $displayName: String!) {\n        tenant: createTenant(tenantID: $id, tenant: {displayName: $displayName}) {\n            id\n            displayName\n        }\n    }\n"): (typeof documents)["\n    mutation CreateTenant($id: String, $displayName: String!) {\n        tenant: createTenant(tenantID: $id, tenant: {displayName: $displayName}) {\n            id\n            displayName\n        }\n    }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    mutation DeleteTenant($id: String!) {\n        id: deleteTenant(tenantID: $id)\n    }\n"): (typeof documents)["\n    mutation DeleteTenant($id: String!) {\n        id: deleteTenant(tenantID: $id)\n    }\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n    mutation UpdateTenant($id: String!, $displayName: String!) {\n        tenant: updateTenant(tenantID: $id, tenant: {displayName: $displayName}) {\n            id\n            displayName\n        }\n    }\n"): (typeof documents)["\n    mutation UpdateTenant($id: String!, $displayName: String!) {\n        tenant: updateTenant(tenantID: $id, tenant: {displayName: $displayName}) {\n            id\n            displayName\n        }\n    }\n"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;