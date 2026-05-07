-- +goose Up
-- +goose StatementBegin

create table system_user
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
insert into system_user (
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
drop table system_user;
-- +goose StatementEnd
