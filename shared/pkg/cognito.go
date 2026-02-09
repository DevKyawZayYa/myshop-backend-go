package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type CognitoConfig struct {
	Region        string
	UserPoolID    string
	AppClientID   string
	RequiredGroup string
}

func (c CognitoConfig) validate() error {
	if strings.TrimSpace(c.Region) == "" {
		return errors.New("region is required")
	}
	if strings.TrimSpace(c.UserPoolID) == "" {
		return errors.New("user pool id is required")
	}
	if strings.TrimSpace(c.AppClientID) == "" {
		return errors.New("app client id is required")
	}
	if strings.TrimSpace(c.RequiredGroup) == "" {
		return errors.New("required group is required")
	}
	return nil
}

func setupCognitoResources(cfg CognitoConfig) (*jwk.Cache, string, string, error) {
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cfg.Region, cfg.UserPoolID)
	jwksURL := fmt.Sprintf("%s/.well-known/jwks.json", issuer)

	cache := jwk.NewCache(context.Background())
	cache.Register(jwksURL)
	if _, err := cache.Refresh(context.Background(), jwksURL); err != nil {
		return nil, "", "", fmt.Errorf("refresh jwks: %w", err)
	}

	return cache, issuer, jwksURL, nil
}

func newCognitoMiddleware(
	cfg CognitoConfig, cache *jwk.Cache, issuer, jwksURL string, validateGroup func(jwt.MapClaims) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := readBearerToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		claims, err := validateAccessToken(c.Request.Context(), accessToken, issuer, jwksURL, cache)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if err := validateClientID(claims, cfg.AppClientID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if validateGroup != nil {
			if err := validateGroup(claims); err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
		}

		c.Next()
	}
}

func NewCognitoJWTMiddleware(cfg CognitoConfig) (gin.HandlerFunc, error) {
	if strings.TrimSpace(cfg.Region) == "" {
		return nil, errors.New("region is required")
	}
	if strings.TrimSpace(cfg.UserPoolID) == "" {
		return nil, errors.New("user pool id is required")
	}
	if strings.TrimSpace(cfg.AppClientID) == "" {
		return nil, errors.New("app client id is required")
	}

	cache, issuer, jwksURL, err := setupCognitoResources(cfg)
	if err != nil {
		return nil, err
	}

	return newCognitoMiddleware(cfg, cache, issuer, jwksURL, nil), nil
}

func NewCognitoAdminMiddleware(cfg CognitoConfig) (gin.HandlerFunc, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cache, issuer, jwksURL, err := setupCognitoResources(cfg)
	if err != nil {
		return nil, err
	}

	groupValidator := func(claims jwt.MapClaims) error {
		return validateGroup(claims, cfg.RequiredGroup)
	}

	return newCognitoMiddleware(cfg, cache, issuer, jwksURL, groupValidator), nil
}

func readBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header must be a bearer token")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("bearer token is required")
	}

	return token, nil
}

func validateAccessToken(ctx context.Context, tokenString, issuer, jwksURL string, cache *jwk.Cache) (jwt.MapClaims, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errors.New("token header missing kid")
		}

		set, err := cache.Get(ctx, jwksURL)
		if err != nil {
			return nil, fmt.Errorf("get jwks: %w", err)
		}

		key, ok := set.LookupKeyID(kid)
		if !ok {
			return nil, errors.New("token key id not found")
		}

		var raw interface{}
		if err := key.Raw(&raw); err != nil {
			return nil, fmt.Errorf("parse jwk: %w", err)
		}

		return raw, nil
	}

	parser := jwt.NewParser(
		jwt.WithIssuer(issuer),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithLeeway(30*time.Second),
	)

	parsed, err := parser.Parse(tokenString, keyfunc)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if tokenUse, ok := claims["token_use"].(string); !ok || tokenUse != "access" {
		return nil, errors.New("token_use must be access")
	}

	return claims, nil
}

func validateClientID(claims jwt.MapClaims, appClientID string) error {
	clientID, ok := claims["client_id"].(string)
	if !ok || clientID == "" {
		return errors.New("client_id claim is missing")
	}
	if clientID != appClientID {
		return errors.New("client_id does not match")
	}
	return nil
}

func validateGroup(claims jwt.MapClaims, requiredGroup string) error {
	groups, ok := claims["cognito:groups"]
	if !ok {
		return errors.New("cognito:groups claim is missing")
	}

	for _, group := range normalizeGroups(groups) {
		if group == requiredGroup {
			return nil
		}
	}

	return errors.New("user is not in required group")
}

func normalizeGroups(groups interface{}) []string {
	switch value := groups.(type) {
	case []string:
		return value
	case []interface{}:
		result := make([]string, 0, len(value))
		for _, item := range value {
			if group, ok := item.(string); ok {
				result = append(result, group)
			}
		}
		return result
	case string:
		return []string{value}
	default:
		return nil
	}
}
