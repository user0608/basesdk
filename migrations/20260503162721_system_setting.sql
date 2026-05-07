-- +goose Up
-- +goose StatementBegin

create table system_properties
(
    key varchar(150) not null primary key,
    value text not null,

    data_type varchar(20) not null
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

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

drop table system_properties;

-- +goose StatementEnd
