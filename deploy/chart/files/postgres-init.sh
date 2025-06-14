#!/bin/bash
# shellcheck disable=SC2154

set -euo pipefail

psql --echo-queries -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_DB}" <<-EOSQL

  -- Create the database & its owner role
  CREATE ROLE greenstar_owner;
  CREATE DATABASE greenstar WITH OWNER=greenstar_owner;
  \c greenstar

  -- Create admin group role
  CREATE ROLE greenstar_admin;
  GRANT ALL PRIVILEGES ON DATABASE greenstar TO greenstar_admin;
  GRANT USAGE ON SCHEMA public TO greenstar_admin;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO greenstar_admin;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO greenstar_admin;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON FUNCTIONS TO greenstar_admin;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TYPES TO greenstar_admin;

  -- Create viewer group role
  CREATE ROLE greenstar_viewer;
  GRANT CONNECT, TEMPORARY ON DATABASE greenstar TO greenstar_viewer;
  GRANT USAGE ON SCHEMA public TO greenstar_viewer;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO greenstar_viewer;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON SEQUENCES TO greenstar_viewer;
  ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE ON TYPES TO greenstar_viewer;

  -- Create greenstar_migration user role
  CREATE ROLE greenstar_migration WITH CREATEROLE LOGIN PASSWORD '{{ .migrationUserPassword }}';
  GRANT CONNECT, TEMPORARY ON DATABASE greenstar TO greenstar_migration WITH GRANT OPTION;
  GRANT CREATE, USAGE ON SCHEMA public TO greenstar_migration WITH GRANT OPTION;

  -- Create the greenstar_server user role
  CREATE ROLE greenstar_server WITH LOGIN PASSWORD '{{ .serverUserPassword }}';
  GRANT CONNECT, TEMPORARY ON DATABASE greenstar TO greenstar_server;
  GRANT USAGE ON SCHEMA public TO greenstar_server;

EOSQL
