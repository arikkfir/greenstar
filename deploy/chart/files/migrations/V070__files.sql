CREATE TABLE files
(
    id           TEXT        NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT        NOT NULL,
    data_oid     oid         NOT NULL,
    size         INTEGER     NOT NULL,
    content_type TEXT        NOT NULL,
    tenant_id    VARCHAR(10) NOT NULL
        CONSTRAINT fk_files_tenant_id
            REFERENCES tenants (id)
            ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT uq_files UNIQUE (tenant_id, name)
);
GRANT SELECT, DELETE, INSERT, UPDATE ON TABLE files TO greenstar_server;
