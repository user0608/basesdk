package usecases_test

import (
	"basesdk/auth/jwt"
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/properties"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/testdb"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type jwtConfigStub struct {
	privateKeyPath string
	publicKeyPath  string
}

var _ configs.JWTConfig = (*jwtConfigStub)(nil)

func (c *jwtConfigStub) GetPrivateKeyPath() string { return c.privateKeyPath }
func (c *jwtConfigStub) GetPublicKeyPath() string  { return c.publicKeyPath }

func TestSecurityUsecaseSystemUserLogin(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	usecase, tokenService := newSecurityUsecase(t, storage)
	ctx := context.Background()

	token, err := usecase.SystemUserLogin(ctx, "kevin", "maira002")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := tokenService.ValidateToken(ctx, token)
	require.NoError(t, err)
	require.Equal(t, jwt.TypeSystem, claims.Type)
	require.Equal(t, "kevin", claims.Subject)
}

func TestSecurityUsecaseTenantUserLogin(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	usecase, tokenService := newSecurityUsecase(t, storage)
	ctx := context.Background()

	token, err := usecase.TenantUserLogin(ctx, "local", "kevin", "maira002")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := tokenService.ValidateToken(ctx, token)
	require.NoError(t, err)
	require.Equal(t, jwt.TypeTenant, claims.Type)
	require.Equal(t, "local", claims.Tenant)
	require.Equal(t, "kevin", claims.Subject)
	require.Equal(t, "America/Lima", claims.TimeZone)
}

func TestSecurityUsecaseTenantUserLoginRejectsInvalidPassword(t *testing.T) {
	storage := testdb.NewPostgresStorage(t)
	usecase, _ := newSecurityUsecase(t, storage)

	token, err := usecase.TenantUserLogin(context.Background(), "local", "kevin", "wrong")
	require.Error(t, err)
	require.Empty(t, token)
}

func newSecurityUsecase(t *testing.T, storage connection.StorageManager) (*usecases.SecurityUsecase, *jwt.TokenService) {
	t.Helper()

	dir := t.TempDir()
	keyStore, err := jwt.NewKeyStore(&jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "private.pem"),
		publicKeyPath:  filepath.Join(dir, "public.pem"),
	})
	require.NoError(t, err)

	systemProps := properties.NewSystemProperties(storage)
	tokenService := jwt.NewTokenService(keyStore, systemProps)
	usecase := usecases.NewSecurityUsecase(
		tokenService,
		repositories.NewSystemUserRepository(storage),
		repositories.NewAppUserRepository(storage),
	)

	return usecase, tokenService
}
