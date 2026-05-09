package handlers_test

import (
	"basesdk/auth/jwt"
	"basesdk/configs"
	"basesdk/connection"
	"basesdk/httpapi"
	"basesdk/properties"
	"basesdk/security/handlers"
	"basesdk/security/repositories"
	"basesdk/security/usecases"
	"basesdk/testdb"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

type loginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func TestSystemLoginE2E(t *testing.T) {
	server, tokenService := newLoginTestServer(t)

	response := executeLoginRequest(t, server, "/api/v1/system/auth/login", `{
		"username": "kevin",
		"password": "maira002"
	}`)

	require.Equal(t, http.StatusOK, response.Code)

	token := decodeLoginToken(t, response)
	require.NotEmpty(t, token)

	claims, err := tokenService.ValidateToken(context.Background(), token)
	require.NoError(t, err)
	require.Equal(t, jwt.TypeSystem, claims.Type)
	require.Equal(t, "kevin", claims.Subject)
}

func TestTenantLoginE2E(t *testing.T) {
	server, tokenService := newLoginTestServer(t)

	response := executeLoginRequest(t, server, "/api/v1/tenants/local/auth/login", `{
		"username": "kevin",
		"password": "maira002"
	}`)

	require.Equal(t, http.StatusOK, response.Code)

	token := decodeLoginToken(t, response)
	require.NotEmpty(t, token)

	claims, err := tokenService.ValidateToken(context.Background(), token)
	require.NoError(t, err)
	require.Equal(t, jwt.TypeTenant, claims.Type)
	require.Equal(t, "local", claims.Tenant)
	require.Equal(t, "kevin", claims.Subject)
	require.Equal(t, "America/Lima", claims.TimeZone)
}

func TestTenantLoginE2ERejectsInvalidPassword(t *testing.T) {
	server, _ := newLoginTestServer(t)

	response := executeLoginRequest(t, server, "/api/v1/tenants/local/auth/login", `{
		"username": "kevin",
		"password": "wrong"
	}`)

	require.GreaterOrEqual(t, response.Code, http.StatusBadRequest)
}

func TestTenantLoginE2ERejectsMissingTenant(t *testing.T) {
	server, _ := newLoginTestServer(t)

	response := executeLoginRequest(t, server, "/api/v1/tenants/missing/auth/login", `{
		"username": "kevin",
		"password": "maira002"
	}`)

	require.GreaterOrEqual(t, response.Code, http.StatusBadRequest)
}

func newLoginTestServer(t *testing.T) (*httptest.Server, *jwt.TokenService) {
	t.Helper()

	storage := testdb.NewPostgresStorage(t)
	usecase, tokenService := newSecurityUsecase(t, storage)

	server := httpapi.NewServer(
		[]httpapi.Route{
			handlers.SystemUserHandler(usecase),
			handlers.TenantUserHandler(usecase),
		},
		nil,
	)

	testServer := httptest.NewServer(server)
	t.Cleanup(testServer.Close)

	return testServer, tokenService
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

func executeLoginRequest(t *testing.T, server *httptest.Server, path string, body string) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(response, request)

	return response
}

func decodeLoginToken(t *testing.T, response *httptest.ResponseRecorder) string {
	t.Helper()

	var payload loginResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))

	return payload.Data.Token
}
