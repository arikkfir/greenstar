#!/usr/bin/env bash

set -exuo pipefail

jq '.exports.import={"types": "./dist/index.d.ts", "default": "./dist/esm/index.js"}' node_modules/@growthbook/growthbook-react/package.json > tmp.json
mv tmp.json node_modules/@growthbook/growthbook-react/package.json
