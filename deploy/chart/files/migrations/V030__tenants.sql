CREATE TABLE tenants
(
    id           TEXT PRIMARY KEY,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    display_name TEXT
);
GRANT SELECT, DELETE, INSERT, UPDATE ON TABLE tenants TO greenstar_server;
