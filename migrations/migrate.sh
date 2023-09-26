#!/usr/bin/env bash

echo "Testing Neo4j connectivity..."
while true; do
  cypher-shell --non-interactive "SHOW DATABASES" >/dev/null && break
  echo "Neo4j is not available, retrying in 6 seconds..."
  sleep 6
done

set -euo pipefail

for f in $(find ./*.cypher | sort); do
  echo "Importing '${f}'..."
  cypher-shell --database=neo4j --fail-fast --format=plain --non-interactive < "${f}"
done

echo "Done!"
