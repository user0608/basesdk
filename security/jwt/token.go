package jwt

import (
	"basesdk/properties"
	"context"
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

const (
	TypeSystem = "system"
	TypeTenant = "tenant"

	defaultTokenTTL       = time.Hour * 720
	defaultIssuer         = "base-sdk"
	defaultAudience       = "app"
	defaultSystemTimeZone = "America/Lima"
)

var (
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrMissingTenant    = errors.New("missing tenant")
	ErrInvalidTimeZone  = errors.New("invalid time zone")
)

type TokenService struct {
	keyStore    *KeyStore
	systemProps *properties.SystemProperties
}

func NewTokenService(
	keyStore *KeyStore,
	systemProps *properties.SystemProperties,
) *TokenService {
	return &TokenService{
		keyStore:    keyStore,
		systemProps: systemProps,
	}
}

type TokenClaims struct {
	Type     string `json:"type"`
	Tenant   string `json:"tenant,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`

	gojwt.RegisteredClaims
}

type TokenData struct {
	Type      string
	Tenant    string
	Username  string
	TimeZone  string
	Issuer    string
	Audience  []string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type tokenConfig struct {
	TTL            time.Duration
	Issuer         string
	Audience       []string
	SystemTimeZone string
}

func (service *TokenService) GenerateSystemToken(ctx context.Context, username string) (string, error) {
	return service.generateToken(ctx, TokenClaims{
		Type: TypeSystem,
		RegisteredClaims: gojwt.RegisteredClaims{
			Subject: username,
		},
	})
}

func (service *TokenService) GenerateTenantToken(ctx context.Context, tenant, username, timeZone string) (string, error) {
	if tenant == "" {
		return "", ErrMissingTenant
	}

	if err := validateTimeZone(timeZone); err != nil {
		return "", err
	}

	return service.generateToken(ctx, TokenClaims{
		Type:     TypeTenant,
		Tenant:   tenant,
		TimeZone: timeZone,
		RegisteredClaims: gojwt.RegisteredClaims{
			Subject: username,
		},
	})
}

func (service *TokenService) ValidateToken(ctx context.Context, rawToken string) (*TokenClaims, error) {
	cfg, err := service.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	claims := &TokenClaims{}

	token, err := gojwt.ParseWithClaims(
		rawToken,
		claims,
		func(token *gojwt.Token) (any, error) {
			if token.Method.Alg() != gojwt.SigningMethodRS512.Alg() {
				return nil, gojwt.ErrTokenSignatureInvalid
			}

			return service.keyStore.VerificationKey(), nil
		},
		gojwt.WithIssuer(cfg.Issuer),
		gojwt.WithAudience(cfg.Audience...),
		gojwt.WithValidMethods([]string{gojwt.SigningMethodRS512.Alg()}),
		gojwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, gojwt.ErrTokenInvalidClaims
	}

	switch claims.Type {
	case TypeSystem:
		if claims.TimeZone == "" {
			claims.TimeZone = cfg.SystemTimeZone
		}

	case TypeTenant:
		if claims.Tenant == "" {
			return nil, ErrMissingTenant
		}

		if claims.TimeZone == "" {
			return nil, ErrInvalidTimeZone
		}

	default:
		return nil, ErrInvalidTokenType
	}

	if err := validateTimeZone(claims.TimeZone); err != nil {
		return nil, err
	}

	return claims, nil
}

func (service *TokenService) GetTokenData(ctx context.Context, rawToken string) (*TokenData, error) {
	claims, err := service.ValidateToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	return claims.Data(), nil
}

func (service *TokenService) RefreshToken(ctx context.Context, rawToken string) (string, error) {
	claims, err := service.ValidateToken(ctx, rawToken)
	if err != nil {
		return "", err
	}

	switch claims.Type {
	case TypeSystem:
		return service.GenerateSystemToken(ctx, claims.Subject)

	case TypeTenant:
		return service.GenerateTenantToken(ctx, claims.Tenant, claims.Subject, claims.TimeZone)

	default:
		return "", ErrInvalidTokenType
	}
}

func (service *TokenService) generateToken(ctx context.Context, claims TokenClaims) (string, error) {
	cfg, err := service.loadConfig(ctx)
	if err != nil {
		return "", err
	}

	switch claims.Type {
	case TypeSystem:
		claims.TimeZone = cfg.SystemTimeZone

	case TypeTenant:
		if claims.Tenant == "" {
			return "", ErrMissingTenant
		}

		if err := validateTimeZone(claims.TimeZone); err != nil {
			return "", err
		}

	default:
		return "", ErrInvalidTokenType
	}

	claims.RegisteredClaims = service.newRegisteredClaims(claims.Subject, cfg)

	token := gojwt.NewWithClaims(gojwt.SigningMethodRS512, claims)

	return token.SignedString(service.keyStore.SigningKey())
}

func (service *TokenService) newRegisteredClaims(username string, cfg tokenConfig) gojwt.RegisteredClaims {
	now := time.Now()

	if cfg.TTL == 0 {
		cfg.TTL = defaultTokenTTL
	}

	if cfg.Issuer == "" {
		cfg.Issuer = defaultIssuer
	}

	if len(cfg.Audience) == 0 {
		cfg.Audience = []string{defaultAudience}
	}

	return gojwt.RegisteredClaims{
		Subject:   username,
		Issuer:    cfg.Issuer,
		Audience:  cfg.Audience,
		IssuedAt:  gojwt.NewNumericDate(now),
		ExpiresAt: gojwt.NewNumericDate(now.Add(cfg.TTL)),
	}
}

func (service *TokenService) loadConfig(ctx context.Context) (tokenConfig, error) {
	ttl, err := service.systemProps.GetDuration(ctx, "jwt_token_ttl", defaultTokenTTL)
	if err != nil {
		return tokenConfig{}, err
	}

	issuer, err := service.systemProps.GetString(ctx, "jwt_issuer", defaultIssuer)
	if err != nil {
		return tokenConfig{}, err
	}

	systemTimeZone, err := service.systemProps.GetString(ctx, "time_zone", defaultSystemTimeZone)
	if err != nil {
		return tokenConfig{}, err
	}

	if systemTimeZone == "" {
		systemTimeZone = defaultSystemTimeZone
	}

	if err := validateTimeZone(systemTimeZone); err != nil {
		return tokenConfig{}, err
	}

	return tokenConfig{
		TTL:            ttl,
		Issuer:         issuer,
		Audience:       []string{defaultAudience},
		SystemTimeZone: systemTimeZone,
	}, nil
}

func validateTimeZone(timeZone string) error {
	if timeZone == "" {
		return ErrInvalidTimeZone
	}

	if _, err := time.LoadLocation(timeZone); err != nil {
		return ErrInvalidTimeZone
	}

	return nil
}

func (claims *TokenClaims) Data() *TokenData {
	data := &TokenData{
		Type:     claims.Type,
		Tenant:   claims.Tenant,
		Username: claims.Subject,
		TimeZone: claims.TimeZone,
		Issuer:   claims.Issuer,
		Audience: claims.Audience,
	}

	if data.TimeZone == "" && claims.Type == TypeSystem {
		data.TimeZone = defaultSystemTimeZone
	}

	if claims.IssuedAt != nil {
		data.IssuedAt = claims.IssuedAt.Time
	}

	if claims.ExpiresAt != nil {
		data.ExpiresAt = claims.ExpiresAt.Time
	}

	return data
}
