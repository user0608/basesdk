package usecases_test

import (
	"basesdk/connection"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/testdb"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthorizationUsecaseAdminPermissionAllowsAll(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	seedAuthorizationAdminUser(t, storage)

	usecase := usecases.NewAuthorizationUsecase(repositories.NewPermissionRepository(storage))
	ctx := context.Background()

	allowed, err := usecase.HasAllPermissions(ctx, "tenant_default", "admin_user", []string{"users.delete", "roles.permissions.replace"})
	require.NoError(t, err)
	require.True(t, allowed)

	allowed, err = usecase.HasAnyPermission(ctx, "tenant_default", "admin_user", []string{"users.delete", "roles.permissions.replace"})
	require.NoError(t, err)
	require.True(t, allowed)
}

func TestAuthorizationUsecaseValidatesUserPermissionSetWithoutAdmin(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	seedAuthorizationUser(t, storage)

	usecase := usecases.NewAuthorizationUsecase(repositories.NewPermissionRepository(storage))
	ctx := context.Background()

	allowed, err := usecase.HasAllPermissions(ctx, "tenant_default", "editor", []string{"users.read", "users.update"})
	require.NoError(t, err)
	require.True(t, allowed)

	allowed, err = usecase.HasAllPermissions(ctx, "tenant_default", "editor", []string{"users.read", "users.delete"})
	require.NoError(t, err)
	require.False(t, allowed)

	allowed, err = usecase.HasAnyPermission(ctx, "tenant_default", "editor", []string{"users.delete", " users.update "})
	require.NoError(t, err)
	require.True(t, allowed)
}

func seedAuthorizationUser(t *testing.T, storage connection.StorageManager) {
	t.Helper()

	err := storage.Conn(context.Background()).Exec(`
		insert into app_user (
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
		values (
			'tenant_default',
			'editor',
			'editor@local',
			'Editor',
			null,
			true,
			false,
			null,
			false,
			'kevin',
			now()
		);

		insert into role (tenant_codigo, code, description, disabled, created_by, created_at)
		values ('tenant_default', 'EDITOR', 'Editor role.', false, 'kevin', now());

		insert into permission (code, description)
		values
			('users.read', 'Read users'),
			('users.update', 'Update users'),
			('users.delete', 'Delete users')
		on conflict (code) do update set description = excluded.description;

		insert into user_role (tenant_codigo, username, role_code, created_by, created_at)
		values ('tenant_default', 'editor', 'EDITOR', 'kevin', now());

		insert into role_permission (tenant_codigo, role_code, permission_code, created_by, created_at)
		values
			('tenant_default', 'EDITOR', 'users.read', 'kevin', now()),
			('tenant_default', 'EDITOR', 'users.update', 'kevin', now());
	`).Error
	require.NoError(t, err)
}

func seedAuthorizationAdminUser(t *testing.T, storage connection.StorageManager) {
	t.Helper()

	err := storage.Conn(context.Background()).Exec(`
		insert into app_user (
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
		values (
			'tenant_default',
			'admin_user',
			'admin_user@local',
			'Admin User',
			null,
			true,
			false,
			null,
			false,
			'kevin',
			now()
		);

		insert into role (tenant_codigo, code, description, disabled, created_by, created_at)
		values ('tenant_default', 'TEST_ADMIN', 'Test admin role.', false, 'kevin', now());

		insert into user_role (tenant_codigo, username, role_code, created_by, created_at)
		values ('tenant_default', 'admin_user', 'TEST_ADMIN', 'kevin', now());

		insert into role_permission (tenant_codigo, role_code, permission_code, created_by, created_at)
		values ('tenant_default', 'TEST_ADMIN', 'admin', 'kevin', now());
	`).Error
	require.NoError(t, err)
}
