-- +goose Up
-- +goose StatementBegin

create table system_properties
(
    key varchar(150) not null primary key,
    value text not null,

    data_type varchar(20) not null,
    constraint system_properties_data_type_check
    check (data_type in ('string', 'int', 'float', 'bool', 'json')),

    description varchar(255)
);

insert into system_properties (
    key,
    value,
    data_type,
    description
) values
(
    'jwt_token_ttl',
    '720h',
    'string',
    'Duración de validez del token JWT (formato duración, ej: 720h)'
);

create table system_account
(
    username varchar(100) not null primary key,
    password_hash varchar(255) not null,

    disabled boolean default false,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone
);

-- usuario por defecto
insert into system_account (
    username,
    password_hash,
    created_by,
    created_at,
    updated_by,
    updated_at
) values (
    'kevin',
    '$2a$12$rpC5L6AQt2x4hUlEqpiwROZQFhOOaAohIzFwPwgElgVDSW01Lcnvu', -- maira002
    'kevin',
    now(),
    'kevin',
    now()
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

drop table if exists system_account;
drop table if exists system_properties;

-- +goose StatementEnd
