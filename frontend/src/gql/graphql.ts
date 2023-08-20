/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
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
  __typename?: 'Account';
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
  __typename?: 'KeyAndValue';
  key: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type KeyAndValueInput = {
  key: Scalars['String']['input'];
  value: Scalars['String']['input'];
};

export type Mutation = {
  __typename?: 'Mutation';
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
  __typename?: 'Operation';
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
  __typename?: 'Query';
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
  __typename?: 'Tenant';
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
  __typename?: 'Transaction';
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

export type RootAccountsQueryVariables = Exact<{
  tenantID: Scalars['ID']['input'];
}>;


export type RootAccountsQuery = { __typename?: 'Query', tenant?: { __typename?: 'Tenant', id: string, accounts: Array<{ __typename?: 'Account', id: string, displayName: string, childCount: number, icon: string }> } | null };

export type AccountChildrenQueryVariables = Exact<{
  tenantID: Scalars['ID']['input'];
  accountID: Scalars['ID']['input'];
}>;


export type AccountChildrenQuery = { __typename?: 'Query', tenant?: { __typename?: 'Tenant', id: string, account?: { __typename?: 'Account', id: string, children: Array<{ __typename?: 'Account', id: string, displayName: string, childCount: number, icon: string }> } | null } | null };

export type AllTenantsQueryVariables = Exact<{ [key: string]: never; }>;


export type AllTenantsQuery = { __typename?: 'Query', tenants: Array<{ __typename?: 'Tenant', id: string, displayName: string }> };

export type CreateTenantMutationVariables = Exact<{
  id?: InputMaybe<Scalars['String']['input']>;
  displayName: Scalars['String']['input'];
}>;


export type CreateTenantMutation = { __typename?: 'Mutation', tenant: { __typename?: 'Tenant', id: string, displayName: string } };

export type DeleteTenantMutationVariables = Exact<{
  id: Scalars['String']['input'];
}>;


export type DeleteTenantMutation = { __typename?: 'Mutation', id: string };

export type UpdateTenantMutationVariables = Exact<{
  id: Scalars['String']['input'];
  displayName: Scalars['String']['input'];
}>;


export type UpdateTenantMutation = { __typename?: 'Mutation', tenant: { __typename?: 'Tenant', id: string, displayName: string } };


export const RootAccountsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"rootAccounts"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"tenantID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"tenant"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"tenantID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"accounts"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"childCount"}},{"kind":"Field","name":{"kind":"Name","value":"icon"}}]}}]}}]}}]} as unknown as DocumentNode<RootAccountsQuery, RootAccountsQueryVariables>;
export const AccountChildrenDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"accountChildren"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"tenantID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"accountID"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"tenant"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"tenantID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"account"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"accountID"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"children"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"childCount"}},{"kind":"Field","name":{"kind":"Name","value":"icon"}}]}}]}}]}}]}}]} as unknown as DocumentNode<AccountChildrenQuery, AccountChildrenQueryVariables>;
export const AllTenantsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"allTenants"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"tenants"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}}]}}]}}]} as unknown as DocumentNode<AllTenantsQuery, AllTenantsQueryVariables>;
export const CreateTenantDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateTenant"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"displayName"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","alias":{"kind":"Name","value":"tenant"},"name":{"kind":"Name","value":"createTenant"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"tenantID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"tenant"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"displayName"},"value":{"kind":"Variable","name":{"kind":"Name","value":"displayName"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}}]}}]}}]} as unknown as DocumentNode<CreateTenantMutation, CreateTenantMutationVariables>;
export const DeleteTenantDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteTenant"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","alias":{"kind":"Name","value":"id"},"name":{"kind":"Name","value":"deleteTenant"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"tenantID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteTenantMutation, DeleteTenantMutationVariables>;
export const UpdateTenantDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateTenant"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"displayName"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","alias":{"kind":"Name","value":"tenant"},"name":{"kind":"Name","value":"updateTenant"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"tenantID"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"tenant"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"displayName"},"value":{"kind":"Variable","name":{"kind":"Name","value":"displayName"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}}]}}]}}]} as unknown as DocumentNode<UpdateTenantMutation, UpdateTenantMutationVariables>;