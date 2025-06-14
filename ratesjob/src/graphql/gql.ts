/* eslint-disable */
import * as types from './graphql.js';
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
 * Learn more about it here: https://the-guild.dev/graphql/codegen/plugins/presets/preset-client#reducing-bundle-size
 */
type Documents = {
    "\n    query GetCurrencies {\n        currencies {\n            code\n        }\n    }\n": typeof types.GetCurrenciesDocument,
    "\n    query GetCurrencyRates($startDate: DateTime, $endDate: DateTime, $sourceCurrencyCode: String, $targetCurrencyCode: String) {\n        currencyRates(\n            startDate: $startDate,\n            endDate: $endDate,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n": typeof types.GetCurrencyRatesDocument,
    "\n    mutation CreateCurrencyRate($date: DateTime!, $sourceCurrencyCode: String!, $targetCurrencyCode: String!, $rate: Float!) {\n        createCurrencyRate(\n            date: $date,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode,\n            rate: $rate\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n": typeof types.CreateCurrencyRateDocument,
};
const documents: Documents = {
    "\n    query GetCurrencies {\n        currencies {\n            code\n        }\n    }\n": types.GetCurrenciesDocument,
    "\n    query GetCurrencyRates($startDate: DateTime, $endDate: DateTime, $sourceCurrencyCode: String, $targetCurrencyCode: String) {\n        currencyRates(\n            startDate: $startDate,\n            endDate: $endDate,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n": types.GetCurrencyRatesDocument,
    "\n    mutation CreateCurrencyRate($date: DateTime!, $sourceCurrencyCode: String!, $targetCurrencyCode: String!, $rate: Float!) {\n        createCurrencyRate(\n            date: $date,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode,\n            rate: $rate\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n": types.CreateCurrencyRateDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    query GetCurrencies {\n        currencies {\n            code\n        }\n    }\n"): (typeof documents)["\n    query GetCurrencies {\n        currencies {\n            code\n        }\n    }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    query GetCurrencyRates($startDate: DateTime, $endDate: DateTime, $sourceCurrencyCode: String, $targetCurrencyCode: String) {\n        currencyRates(\n            startDate: $startDate,\n            endDate: $endDate,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n"): (typeof documents)["\n    query GetCurrencyRates($startDate: DateTime, $endDate: DateTime, $sourceCurrencyCode: String, $targetCurrencyCode: String) {\n        currencyRates(\n            startDate: $startDate,\n            endDate: $endDate,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    mutation CreateCurrencyRate($date: DateTime!, $sourceCurrencyCode: String!, $targetCurrencyCode: String!, $rate: Float!) {\n        createCurrencyRate(\n            date: $date,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode,\n            rate: $rate\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n"): (typeof documents)["\n    mutation CreateCurrencyRate($date: DateTime!, $sourceCurrencyCode: String!, $targetCurrencyCode: String!, $rate: Float!) {\n        createCurrencyRate(\n            date: $date,\n            sourceCurrencyCode: $sourceCurrencyCode,\n            targetCurrencyCode: $targetCurrencyCode,\n            rate: $rate\n        ) {\n            date\n            sourceCurrency {\n                code\n            }\n            targetCurrency {\n                code\n            }\n            rate\n        }\n    }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;