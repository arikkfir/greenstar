import type { CodegenConfig } from "@graphql-codegen/cli"

const config: CodegenConfig = {
    schema: "../server/src/schema/*.graphql",
    documents: ["src/**/*.ts"],
    ignoreNoDocuments: true,
    emitLegacyCommonJSImports: false,
    config: {
        importModuleSpecifierEnding: 'js',
    },
    generates: {
        "./src/graphql/": {
            preset: "client",
            plugins: [
                {
                    add: {
                        content: 'import { DateTime } from "luxon";',
                    },
                },
            ],
            config: {
                scalars: {
                    DateTime: "DateTime",
                },
            },
            presetConfig: {
                gqlTagName: "gql",
                scalars: {
                    DateTime: "DateTime",
                },
            },
        },
    },
}

export default config
