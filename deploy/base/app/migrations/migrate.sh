#!/usr/bin/env bash

set -euxo pipefail

CMD=$(which cypher-shell)

if [[ -z "${CMD}" ]]; then
  echo "could not find cypher-shell" >&2 && exit 1
fi

pwd
ls -la

for f in $(ls *.cypher | sort); do
  ${CMD}  --address=neo4j://greenstar-neo4j \
          --database=neo4j \
          --fail-fast \
          --format=plain \
          --non-interactive \
          < ${f}
done
