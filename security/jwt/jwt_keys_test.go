package jwt

import (
	"basesdk/configs"
	"os"
	"path/filepath"

	"testing"
)

type jwtConfigStub struct {
	privateKeyPath string
	publicKeyPath  string
}

var _ configs.JWTConfig = (*jwtConfigStub)(nil)

func (c jwtConfigStub) GetPrivateKeyPath() string { return c.privateKeyPath }
func (c jwtConfigStub) GetPublicKeyPath() string  { return c.publicKeyPath }
func (c jwtConfigStub) GetJwtTokenTTL() string    { return "" }
func (c jwtConfigStub) GetIssuer() string         { return "" }

func TestNewJWTKeysStoreGeneratesKeysWhenMissing(t *testing.T) {
	dir := t.TempDir()

	config := jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "jwt", "private.pem"),
		publicKeyPath:  filepath.Join(dir, "jwt", "public.pem"),
	}

	provider, err := NewKeyStore(config)
	if err != nil {
		t.Fatalf("NewKeyStore() error = %v", err)
	}

	if provider.SigningKey() == nil {
		t.Fatal("SigningKey() returned nil")
	}

	if provider.SigningKey() == nil {
		t.Fatal("VerificationKey() returned nil")
	}

	if _, err := os.Stat(config.privateKeyPath); err != nil {
		t.Fatalf("private key file was not created: %v", err)
	}

	if _, err := os.Stat(config.publicKeyPath); err != nil {
		t.Fatalf("public key file was not created: %v", err)
	}
}

func TestNewJWTKeysStoreLoadsExistingKeys(t *testing.T) {
	dir := t.TempDir()

	config := jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "private.pem"),
		publicKeyPath:  filepath.Join(dir, "public.pem"),
	}

	firstProvider, err := NewKeyStore(config)
	if err != nil {
		t.Fatalf("first NewKeyStore() error = %v", err)
	}

	secondProvider, err := NewKeyStore(config)
	if err != nil {
		t.Fatalf("second NewKeyStore() error = %v", err)
	}

	if firstProvider.SigningKey().N.Cmp(secondProvider.SigningKey().N) != 0 {
		t.Fatal("expected existing private key to be loaded, got different modulus")
	}

	if firstProvider.VerificationKey().N.Cmp(secondProvider.VerificationKey().N) != 0 {
		t.Fatal("expected existing public key to be loaded, got different modulus")
	}
}

func TestNewJWTKeysStoreRegeneratesPairWhenOnlyPrivateExists(t *testing.T) {
	dir := t.TempDir()

	config := jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "private.pem"),
		publicKeyPath:  filepath.Join(dir, "public.pem"),
	}

	if err := os.MkdirAll(filepath.Dir(config.privateKeyPath), 0700); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(config.privateKeyPath, []byte("invalid private key"), 0600); err != nil {
		t.Fatal(err)
	}

	provider, err := NewKeyStore(config)
	if err != nil {
		t.Fatalf("NewKeyStore() error = %v", err)
	}

	if provider.SigningKey() == nil {
		t.Fatal("SigningKey() returned nil")
	}

	if provider.VerificationKey() == nil {
		t.Fatal("VerificationKey() returned nil")
	}
}

func TestNewJWTKeysStoreRegeneratesPairWhenOnlyPublicExists(t *testing.T) {
	dir := t.TempDir()

	config := jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "private.pem"),
		publicKeyPath:  filepath.Join(dir, "public.pem"),
	}

	if err := os.MkdirAll(filepath.Dir(config.publicKeyPath), 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(config.publicKeyPath, []byte("invalid public key"), 0644); err != nil {
		t.Fatal(err)
	}

	provider, err := NewKeyStore(config)
	if err != nil {
		t.Fatalf("NewKeyStore() error = %v", err)
	}

	if provider.SigningKey() == nil {
		t.Fatal("SigningKey() returned nil")
	}

	if provider.VerificationKey() == nil {
		t.Fatal("VerificationKey() returned nil")
	}
}

func TestNewJWTKeysStoreReturnsErrorWhenExistingPrivateKeyIsInvalid(t *testing.T) {
	dir := t.TempDir()

	config := jwtConfigStub{
		privateKeyPath: filepath.Join(dir, "private.pem"),
		publicKeyPath:  filepath.Join(dir, "public.pem"),
	}

	if err := os.WriteFile(config.privateKeyPath, []byte("invalid private key"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(config.publicKeyPath, []byte("invalid public key"), 0644); err != nil {
		t.Fatal(err)
	}

	provider, err := NewKeyStore(config)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if provider != nil {
		t.Fatal("expected provider nil")
	}
}
