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


create table tenant_system_properties
(
    key varchar(150) not null,
    value text not null,
    tenant_codigo varchar(100) not null,
    data_type varchar(20) not null
    check (data_type in ('string', 'int', 'float', 'bool', 'json')),

    description varchar(255),

    primary key (tenant_codigo, key),
    constraint fk_tenant_system_properties_tenant
    foreign key (tenant_codigo) references tenant (codigo)
);

create table app_user
(
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

    constraint pk_app_user
    primary key (tenant_codigo, username),

    constraint fk_app_user_tenant
    foreign key (tenant_codigo)
    references tenant (codigo),

    constraint uq_app_user_tenant_email
    unique (tenant_codigo, email)
);

comment on table app_user is 'User account that belongs to a tenant.';
comment on column app_user.username is 'Unique username inside the tenant.';
comment on column app_user.password_hash is 'Hashed password. Never store plain text passwords.';
comment on column app_user.email_verified is 'Indicates whether the user email has been verified.';
comment on column app_user.must_change_password is 'Forces the user to change password on next login.';
comment on column app_user.last_login_at is 'Timestamp of the latest successful login.';
comment on column app_user.disabled is 'Disables user access without deleting the user.';

create table role
(
    tenant_codigo varchar(100) not null,
    code varchar(100) not null,

    description varchar(500),
    disabled boolean not null default false,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint pk_role
    primary key (tenant_codigo, code),

    constraint fk_role_tenant
    foreign key (tenant_codigo)
    references tenant (codigo)
);

comment on table role is 'Named group of permissions assigned to users within a tenant.';
comment on column role.code is 'Role code unique inside the tenant, for example SUPER_ADMIN.';
comment on column role.disabled is 'Disables the role without deleting its permissions or assignments.';

create table permission
(
    code varchar(150) not null,

    description varchar(500),

    constraint pk_permission
    primary key (code)
);

comment on table permission is 'Global permission catalog used by roles.';
comment on column permission.code is 'Permission code, for example users.create, users.read, roles.update.';

create table user_role
(
    tenant_codigo varchar(100) not null,
    username varchar(100) not null,
    role_code varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint pk_user_role
    primary key (tenant_codigo, username, role_code),

    constraint fk_user_role_user
    foreign key (tenant_codigo, username)
    references app_user (tenant_codigo, username),

    constraint fk_user_role_role
    foreign key (tenant_codigo, role_code)
    references role (tenant_codigo, code)
);

comment on table user_role is 'Assigns roles directly to users inside a tenant.';

create table role_permission
(
    tenant_codigo varchar(100) not null,
    role_code varchar(100) not null,
    permission_code varchar(150) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint pk_role_permission
    primary key (tenant_codigo, role_code, permission_code),

    constraint fk_role_permission_role
    foreign key (tenant_codigo, role_code)
    references role (tenant_codigo, code),

    constraint fk_role_permission_permission
    foreign key (permission_code)
    references permission (code)
);

comment on table role_permission is 'Assigns permissions to roles inside a tenant.';

create table app_group
(
    tenant_codigo varchar(100) not null,
    code varchar(100) not null,

    description varchar(500),
    disabled boolean not null default false,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint pk_app_group
    primary key (tenant_codigo, code),

    constraint fk_app_group_tenant
    foreign key (tenant_codigo)
    references tenant (codigo)
);

comment on table app_group is 'Tenant group used to organize users by team, area, department, branch, or business unit.';
comment on column app_group.code is 'Group code unique inside the tenant.';
comment on column app_group.disabled is 'Disables the group without deleting its users or role assignments.';

create table user_group
(
    tenant_codigo varchar(100) not null,
    username varchar(100) not null,
    group_code varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint pk_user_group
    primary key (tenant_codigo, username, group_code),

    constraint fk_user_group_user
    foreign key (tenant_codigo, username)
    references app_user (tenant_codigo, username),

    constraint fk_user_group_group
    foreign key (tenant_codigo, group_code)
    references app_group (tenant_codigo, code)
);

comment on table user_group is 'Assigns users to tenant groups.';

create table group_role
(
    tenant_codigo varchar(100) not null,
    group_code varchar(100) not null,
    role_code varchar(100) not null,

    created_by varchar(100) not null,
    created_at timestamp without time zone not null,
    updated_by varchar(100),
    updated_at timestamp without time zone,

    constraint pk_group_role
    primary key (tenant_codigo, group_code, role_code),

    constraint fk_group_role_group
    foreign key (tenant_codigo, group_code)
    references app_group (tenant_codigo, code),

    constraint fk_group_role_role
    foreign key (tenant_codigo, role_code)
    references role (tenant_codigo, code)
);

comment on table group_role is 'Assigns roles to groups so users can inherit permissions through group membership.';

insert into app_user
(
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
    tenant_codigo,
    code,
    description,
    disabled,
    created_by,
    created_at
)
values
(
    'tenant_default',
    'SUPER_ADMIN',
    'Full access role.',
    false,
    'kevin',
    now()
);

insert into user_role
(
    tenant_codigo,
    username,
    role_code,
    created_by,
    created_at
)
values
(
    'tenant_default',
    'kevin',
    'SUPER_ADMIN',
    'kevin',
    now()
);

insert into role_permission
(
    tenant_codigo,
    role_code,
    permission_code,
    created_by,
    created_at
)
select
    role_super_admin.tenant_codigo,
    role_super_admin.code as role_code,
    permission.code as permission_code,
    'kevin' as created_by,
    now() as created_at
from role role_super_admin
cross join permission
where
    role_super_admin.tenant_codigo = 'tenant_default'
    and role_super_admin.code = 'SUPER_ADMIN';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists tenant_system_properties;
drop table if exists group_role;
drop table if exists user_group;
drop table if exists app_group;
drop table if exists role_permission;
drop table if exists user_role;
drop table if exists permission;
drop table if exists role;
drop table if exists app_user;
drop table if exists tenant;

-- +goose StatementEnd
