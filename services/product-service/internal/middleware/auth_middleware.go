package middleware

import (
	"errors"
	"os"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

func RegisterAuthMiddleware() (gin.HandlerFunc, error) {
	cfg := pkg.CognitoConfig{
		Region:        os.Getenv("COGNITO_REGION"),
		UserPoolID:    os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:   os.Getenv("COGNITO_APP_CLIENT_ID"),
		RequiredGroup: "admin",
	}

	if cfg.Region == "" || cfg.UserPoolID == "" || cfg.AppClientID == "" {
		return nil, errors.New("missing required cognito env vars")
	}

	return pkg.NewCognitoAdminMiddleware(cfg)
}
