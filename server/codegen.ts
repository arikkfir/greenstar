import type { CodegenConfig } from "@graphql-codegen/cli"

const config: CodegenConfig = {
    overwrite: true,
    schema: "./src/schema/*.graphql",
    require: ["luxon"],
    generates: {
        "src/schema/graphql.ts": {
            plugins: [
                "typescript",
                "typescript-resolvers",
                "typescript-document-nodes",
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
        },
    },
}

export default config
