import {GraphQLClient} from 'graphql-request';
import {GraphQLClientRequestHeaders} from 'graphql-request';
import gql from 'graphql-tag';

export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
    ID: { input: string; output: string; }
    String: { input: string; output: string; }
    Boolean: { input: boolean; output: boolean; }
    Int: { input: number; output: number; }
    Float: { input: number; output: number; }
    DateTime: { input: any; output: any; }
    Money: { input: any; output: any; }
    Time: { input: any; output: any; }
};

export type Account = {
    childCount: Scalars['Int']['output'];
    children: Array<Account>;
    displayName: Scalars['String']['output'];
    icon: Scalars['String']['output'];
    id: Scalars['ID']['output'];
    incomingTransactions: Array<Transaction>;
    labels: Array<KeyAndValue>;
    outgoingTransactions: Array<Transaction>;
    parent?: Maybe<Account>;
    tenant: Tenant;
};

export type AccountChanges = {
    displayName?: InputMaybe<Scalars['String']['input']>;
    icon?: InputMaybe<Scalars['String']['input']>;
    labels?: InputMaybe<Array<KeyAndValueInput>>;
    parentID?: InputMaybe<Scalars['ID']['input']>;
};

export type KeyAndValue = {
    key: Scalars['String']['output'];
    value: Scalars['String']['output'];
};

export type KeyAndValueInput = {
    key: Scalars['String']['input'];
    value: Scalars['String']['input'];
};

export type Mutation = {
    createAccount: Account;
    createTenant: Tenant;
    createTransaction: Transaction;
    createTransactions: Scalars['Int']['output'];
    deleteAccount: Scalars['ID']['output'];
    deleteTenant: Scalars['ID']['output'];
    scrapeIsraelBankYahav: Scalars['String']['output'];
    updateAccount: Account;
    updateOperation: Operation;
    updateTenant: Tenant;
};


export type MutationCreateAccountArgs = {
    account: AccountChanges;
    accountID?: InputMaybe<Scalars['ID']['input']>;
    tenantID: Scalars['ID']['input'];
};


export type MutationCreateTenantArgs = {
    tenant: TenantChanges;
    tenantID?: InputMaybe<Scalars['String']['input']>;
};


export type MutationCreateTransactionArgs = {
    tenantID: Scalars['ID']['input'];
    transaction: TransactionChanges;
};


export type MutationCreateTransactionsArgs = {
    tenantID: Scalars['ID']['input'];
    transactions: Array<TransactionChanges>;
};


export type MutationDeleteAccountArgs = {
    accountID: Scalars['ID']['input'];
    tenantID: Scalars['ID']['input'];
};


export type MutationDeleteTenantArgs = {
    tenantID: Scalars['String']['input'];
};


export type MutationScrapeIsraelBankYahavArgs = {
    id: Scalars['String']['input'];
    password: Scalars['String']['input'];
    tenantID: Scalars['ID']['input'];
    username: Scalars['String']['input'];
};


export type MutationUpdateAccountArgs = {
    account: AccountChanges;
    accountID: Scalars['ID']['input'];
    tenantID: Scalars['ID']['input'];
};


export type MutationUpdateOperationArgs = {
    id: Scalars['ID']['input'];
    op: OperationChanges;
};


export type MutationUpdateTenantArgs = {
    tenant: TenantChanges;
    tenantID: Scalars['String']['input'];
};

export type Operation = {
    createdAt: Scalars['DateTime']['output'];
    description?: Maybe<Scalars['String']['output']>;
    id: Scalars['ID']['output'];
    name: Scalars['String']['output'];
    result: OperationResult;
    status: OperationStatus;
    updatedAt: Scalars['DateTime']['output'];
};

export type OperationChanges = {
    description?: InputMaybe<Scalars['String']['input']>;
    name: Scalars['String']['input'];
    result: OperationResult;
    status: OperationStatus;
};

export enum OperationResult {
    Failed = 'FAILED',
    Succeeded = 'SUCCEEDED'
}

export enum OperationStatus {
    Accepted = 'ACCEPTED',
    Completed = 'COMPLETED',
    Pending = 'PENDING',
    Rejected = 'REJECTED',
    Started = 'STARTED'
}

export type Query = {
    operation?: Maybe<Operation>;
    tenant?: Maybe<Tenant>;
    tenants: Array<Tenant>;
};


export type QueryOperationArgs = {
    id: Scalars['ID']['input'];
};


export type QueryTenantArgs = {
    id: Scalars['ID']['input'];
};

export type Tenant = {
    account?: Maybe<Account>;
    accounts: Array<Account>;
    displayName: Scalars['String']['output'];
    id: Scalars['ID']['output'];
    transactions: Array<Transaction>;
};


export type TenantAccountArgs = {
    id: Scalars['ID']['input'];
};

export type TenantChanges = {
    displayName: Scalars['String']['input'];
};

export type Transaction = {
    Date: Scalars['Time']['output'];
    amount: Scalars['Money']['output'];
    description: Scalars['String']['output'];
    id: Scalars['ID']['output'];
    referenceID: Scalars['String']['output'];
    sourceAccount: Account;
    targetAccount: Account;
};

export type TransactionChanges = {
    Date: Scalars['Time']['input'];
    amount: Scalars['Money']['input'];
    description: Scalars['String']['input'];
    referenceID: Scalars['String']['input'];
    sourceAccountID: Scalars['ID']['input'];
    targetAccountID: Scalars['ID']['input'];
};

export type LoadTenantQueryVariables = Exact<{
    tenantID: Scalars['ID']['input'];
}>;


export type LoadTenantQuery = { tenant?: { id: string, displayName: string } | null };

export type CreateTenantMutationVariables = Exact<{
    tenantID?: InputMaybe<Scalars['String']['input']>;
    displayName: Scalars['String']['input'];
}>;


export type CreateTenantMutation = { createTenant: { id: string, displayName: string } };

export type DeleteTenantMutationVariables = Exact<{
    tenantID: Scalars['String']['input'];
}>;


export type DeleteTenantMutation = { deleteTenant: string };


export const LoadTenantDocument = gql`
    query loadTenant($tenantID: ID!) {
        tenant(id: $tenantID) {
            id
            displayName
        }
    }
`;
export const CreateTenantDocument = gql`
    mutation createTenant($tenantID: String, $displayName: String!) {
        createTenant(tenantID: $tenantID, tenant: {displayName: $displayName}) {
            id
            displayName
        }
    }
`;
export const DeleteTenantDocument = gql`
    mutation deleteTenant($tenantID: String!) {
        deleteTenant(tenantID: $tenantID)
    }
`;

export type SdkFunctionWrapper = <T>(action: (requestHeaders?: Record<string, string>) => Promise<T>, operationName: string, operationType?: string) => Promise<T>;


const defaultWrapper: SdkFunctionWrapper = (action, _operationName, _operationType) => action();

export function getSdk(client: GraphQLClient, withWrapper: SdkFunctionWrapper = defaultWrapper) {
    return {
        loadTenant(variables: LoadTenantQueryVariables, requestHeaders?: GraphQLClientRequestHeaders): Promise<LoadTenantQuery> {
            return withWrapper((wrappedRequestHeaders) => client.request<LoadTenantQuery>(LoadTenantDocument, variables, {...requestHeaders, ...wrappedRequestHeaders}), 'loadTenant', 'query');
        },
        createTenant(variables: CreateTenantMutationVariables, requestHeaders?: GraphQLClientRequestHeaders): Promise<CreateTenantMutation> {
            return withWrapper((wrappedRequestHeaders) => client.request<CreateTenantMutation>(CreateTenantDocument, variables, {...requestHeaders, ...wrappedRequestHeaders}), 'createTenant', 'mutation');
        },
        deleteTenant(variables: DeleteTenantMutationVariables, requestHeaders?: GraphQLClientRequestHeaders): Promise<DeleteTenantMutation> {
            return withWrapper((wrappedRequestHeaders) => client.request<DeleteTenantMutation>(DeleteTenantDocument, variables, {...requestHeaders, ...wrappedRequestHeaders}), 'deleteTenant', 'mutation');
        }
    };
}

export type Sdk = ReturnType<typeof getSdk>;