#!/usr/bin/env bash

[[ -z "${NEO4J_ADDRESS}" ]] && echo "NEO4J_ADDRESS is not set" && exit 1

echo "Testing Neo4j connectivity..."
while true; do
  cypher-shell --address="${NEO4J_ADDRESS}" --non-interactive "SHOW DATABASES" >/dev/null && break
  echo "Neo4j is not available, retrying in 6 seconds..."
  sleep 6
done

set -euo pipefail

for f in $(find ./*.cypher | sort); do
  echo "Importing '${f}'..."
  cypher-shell  --address="${NEO4J_ADDRESS}" \
                --database=neo4j \
                --fail-fast \
                --format=plain \
                --non-interactive \
                < "${f}"
done

echo "Done!"
