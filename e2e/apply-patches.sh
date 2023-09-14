#!/usr/bin/env bash

set -euo pipefail

patch node_modules/graphql-request/build/esm/types.d.ts <<EOF
5a6,7
> import * as https from "https";
> import * as http from "http";
EOF

patch node_modules/graphql-request/build/esm/types.d.ts <<EOF
60a61,61
>     agent?: https.Agent | http.Agent
EOF

patch node_modules/graphql-request/src/types.ts <<EOF
5a6,7
> import * as https from "https";
> import * as http from "http";
EOF

patch node_modules/graphql-request/src/types.ts <<EOF
96a97,97
>   agent?: https.Agent | http.Agent
EOF

patch node_modules/graphql-request/build/esm/index.d.ts <<EOF
141a142,142
> export type GraphQLClientRequestHeaders = Headers | string[][] | Record<string, string>;
EOF
