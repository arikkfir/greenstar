create table currencies
(
    code           varchar(3)   not null primary key,
    created_at     timestamp with time zone default now(),
    updated_at     timestamp with time zone default now(),
    symbol         varchar(3)   not null,
    symbol_native  text         not null,
    name           text         not null,
    decimal_digits int          not null,
    rounding       int          not null,
    name_plural    text         not null,
    type           text         not null,
    country_codes  varchar(2)[] not null
);
ALTER TABLE currencies
    OWNER TO greenstar_backend;

create table rates
(
    date                 date    not null,
    source_currency_code text    not null
        constraint fk_rates_source_currency_code
            references currencies
            on update cascade on delete cascade,
    target_currency_code text    not null
        constraint fk_rates_target_currency_code
            references currencies
            on update cascade on delete cascade,
    created_at           timestamp with time zone default now(),
    updated_at           timestamp with time zone default now(),
    rate                 numeric not null,
    mock                 boolean not null         default true,
    line_number          int,
    primary key (date, source_currency_code, target_currency_code)
);
ALTER TABLE rates
    OWNER TO greenstar_backend;

create table tenants
(
    id           text primary key,
    created_at   timestamp with time zone default now(),
    updated_at   timestamp with time zone default now(),
    display_name text
);
ALTER TABLE tenants
    OWNER TO greenstar_backend;

create table accounts
(
    id           uuid primary key         default gen_random_uuid(),
    created_at   timestamp with time zone default now(),
    updated_at   timestamp with time zone default now(),
    display_name text not null check (char_length(display_name) > 0),
    icon         text check (char_length(icon) > 0),
    parent_id    uuid
        constraint fk_accounts_parent_id
            references accounts
            on update cascade on delete cascade,
    tenant_id    text
        constraint fk_accounts_tenant_id
            references tenants
            on update cascade on delete cascade
);
ALTER TABLE accounts
    OWNER TO greenstar_backend;

create table transactions
(
    id                uuid primary key         default gen_random_uuid(),
    created_at        timestamp with time zone default now(),
    updated_at        timestamp with time zone default now(),
    date              timestamp with time zone,
    reference_id      text       not null check (char_length(reference_id) > 0),
    amount            numeric    not null check (amount > 0),
    currency          varchar(3) not null check (char_length(currency) = 3)
        constraint fk_transactions_currency
            references currencies
            on update cascade on delete restrict,
    description       text CHECK (char_length(description) > 0),
    source_account_id uuid       not null
        constraint fk_transactions_source_account_id
            references accounts
            on update cascade on delete cascade,
    target_account_id uuid       not null
        constraint fk_transactions_target_account_id
            references accounts
            on update cascade on delete cascade
);
ALTER TABLE transactions
    OWNER TO greenstar_backend;
