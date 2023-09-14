#!/usr/bin/env bash

set -exuo pipefail

(cd backend && go generate ./...)
(cd frontend && npm run codegen)
(cd e2e && npm run codegen)
