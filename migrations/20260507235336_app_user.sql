-- +goose Up
-- +goose StatementBegin

create table app_user
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    username varchar(100) not null,
    email varchar(255) not null,
    full_name varchar(255),

    password_hash varchar(255),
    email_verified boolean not null default false,
    must_change_password boolean not null default false,
    last_login_at timestamp without time zone,

    disabled boolean not null default false,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_app_user_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint uq_app_user_tenant_username
    unique (tenant_codigo, username),

    constraint uq_app_user_tenant_email
    unique (tenant_codigo, email)
);

comment on table app_user is 'User account that belongs to a tenant.';
comment on column app_user.password_hash is 'Hashed password. Never store plain text passwords.';
comment on column app_user.email_verified is 'Indicates whether the user email has been verified.';
comment on column app_user.must_change_password is 'Forces the user to change password on next login.';
comment on column app_user.last_login_at is 'Timestamp of the latest successful login.';
comment on column app_user.disabled is 'Disables user access without deleting the user.';

create table role
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    name varchar(100) not null,
    description varchar(500),
    disabled boolean not null default false,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_role_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint uq_role_tenant_name
    unique (tenant_codigo, name)
);

comment on table role is 'Named group of permissions assigned to users within a tenant.';
comment on column role.disabled is 'Disables the role without deleting its permissions or assignments.';

create table permission
(
    codigo varchar(100) not null primary key,

    name varchar(150) not null,
    description varchar(500),

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint uq_permission_name
    unique (name)
);

comment on table permission is 'Global permission catalog used by roles.';
comment on column permission.name is 'Permission key, for example users.create, users.read, roles.update.';

create table user_role
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    user_codigo varchar(100) not null,
    role_codigo varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_user_role_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint fk_user_role_user
    foreign key (user_codigo)
    references app_user (codigo),

    constraint fk_user_role_role
    foreign key (role_codigo)
    references role (codigo),

    constraint uq_user_role
    unique (tenant_codigo, user_codigo, role_codigo)
);

comment on table user_role is 'Assigns roles directly to users inside a tenant.';

create table role_permission
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    role_codigo varchar(100) not null,
    permission_codigo varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_role_permission_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint fk_role_permission_role
    foreign key (role_codigo)
    references role (codigo),

    constraint fk_role_permission_permission
    foreign key (permission_codigo)
    references permission (codigo),

    constraint uq_role_permission
    unique (tenant_codigo, role_codigo, permission_codigo)
);

comment on table role_permission is 'Assigns permissions to roles inside a tenant.';

create table app_group
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    name varchar(100) not null,
    description varchar(500),
    disabled boolean not null default false,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_app_group_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint uq_app_group_tenant_name
    unique (tenant_codigo, name)
);

comment on table app_group is 'Tenant group used to organize users by team, area, department, branch, or business unit.';
comment on column app_group.disabled is 'Disables the group without deleting its users or role assignments.';

create table user_group
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    user_codigo varchar(100) not null,
    group_codigo varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_user_group_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint fk_user_group_user
    foreign key (user_codigo)
    references app_user (codigo),

    constraint fk_user_group_group
    foreign key (group_codigo)
    references app_group (codigo),

    constraint uq_user_group
    unique (tenant_codigo, user_codigo, group_codigo)
);

comment on table user_group is 'Assigns users to tenant groups.';

create table group_role
(
    codigo varchar(100) not null primary key,
    tenant_codigo varchar(100) not null,

    group_codigo varchar(100) not null,
    role_codigo varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint fk_group_role_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint fk_group_role_group
    foreign key (group_codigo)
    references app_group (codigo),

    constraint fk_group_role_role
    foreign key (role_codigo)
    references role (codigo),

    constraint uq_group_role
    unique (tenant_codigo, group_codigo, role_codigo)
);

comment on table group_role is 'Assigns roles to groups so users can inherit permissions through group membership.';

insert into app_user
(
    codigo,
    tenant_codigo,
    username,
    email,
    full_name,
    password_hash,
    email_verified,
    must_change_password,
    last_login_at,
    disabled,
    created_by,
    created_at
)
values
(
    'user_kevin',
    'tenant_default',
    'kevin',
    'kevin@local',
    'Kevin',
    '$2a$12$rpC5L6AQt2x4hUlEqpiwROZQFhOOaAohIzFwPwgElgVDSW01Lcnvu',
    true,
    false,
    null,
    false,
    'kevin',
    now()
);

insert into role
(
    codigo,
    tenant_codigo,
    name,
    description,
    disabled,
    created_by,
    created_at
)
values
(
    'role_super_admin',
    'tenant_default',
    'SUPER_ADMIN',
    'Full access role.',
    false,
    'kevin',
    now()
);

insert into user_role
(
    codigo,
    tenant_codigo,
    user_codigo,
    role_codigo,
    created_by,
    created_at
)
values
(
    'user_role_kevin_super_admin',
    'tenant_default',
    'user_kevin',
    'role_super_admin',
    'kevin',
    now()
);

insert into role_permission
(
    codigo,
    tenant_codigo,
    role_codigo,
    permission_codigo,
    created_by,
    created_at
)
select
    'rp_' || role_super_admin.codigo || '_' || permission.codigo as codigo,
    role_super_admin.tenant_codigo,
    role_super_admin.codigo as role_codigo,
    permission.codigo as permission_codigo,
    'kevin' as created_by,
    now() as created_at
from role role_super_admin
cross join permission
where role_super_admin.codigo = 'role_super_admin';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists group_role;
drop table if exists user_group;
drop table if exists app_group;
drop table if exists role_permission;
drop table if exists user_role;
drop table if exists permission;
drop table if exists role;
drop table if exists app_user;

-- +goose StatementEnd
