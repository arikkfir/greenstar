import type {CodegenConfig} from '@graphql-codegen/cli';

const config: CodegenConfig = {
    overwrite: true,
    schema: "../schema",
    documents: [
        "src/**/*.graphql",
        "src/**/*.ts",
    ],
    generates: {
        "src/generated/graphql.ts": {
            plugins: [
                {
                    "typescript": {
                        emitLegacyCommonJSImports: false,
                        skipTypename: true,
                    },
                },
                {
                    "typescript-operations": {
                        emitLegacyCommonJSImports: false,
                        skipTypename: true,
                    },
                },
                {
                    "typescript-graphql-request": {
                        emitLegacyCommonJSImports: false,
                        skipTypename: true,
                    },

                }
            ],
        },
    },
};

export default config;
