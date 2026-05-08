-- +goose Up
-- +goose StatementBegin

create table tenant
(
    codigo varchar(100) not null primary key,
    name varchar(255) not null,
    timezone varchar(100) not null,
    max_active_users integer not null,
    disabled boolean not null default false,
    expires_at timestamp without time zone,
    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone
);

comment on table tenant is 'Represents an organization, company, or isolated customer account in the multi-tenant system.';
comment on column tenant.timezone is 'Tenant timezone used to display dates, schedule jobs, and calculate local business time.';
comment on column tenant.max_active_users is 'Maximum number of active users allowed for this tenant.';
comment on column tenant.disabled is 'Disables access to the tenant without deleting its data.';
comment on column tenant.expires_at is 'Expiration date for the tenant license, trial, or subscription.';

insert into tenant
(
    codigo,
    name,
    timezone,
    max_active_users,
    disabled,
    expires_at,
    created_by,
    created_at
)
values
(
    'tenant_default',
    'Default Tenant',
    'America/Lima',
    999999,
    false,
    null,
    'kevin',
    now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists tenant;

-- +goose StatementEnd
