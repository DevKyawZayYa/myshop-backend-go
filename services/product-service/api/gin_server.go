package server

import (
	"fmt"
	"log"
	"os"
	"product-service/config"
	router "product-service/internal"
	"product-service/internal/middleware"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	db := config.ConnectDB()

	db.AutoMigrate(
	// &entity.Category{},
	// &entity.Product{},
	// &entity.Attribute{},
	// &entity.AttributeValue{},
	// &entity.Variant{},
	// &entity.ProductImage{},
	)

	r := gin.Default()

	// Configure validator to use JSON field names in error messages
	pkg.ConfigureValidator(r)

	// Register middleware
	middleware.RegisterBasicMiddleware(r)

	jwtAuth, err := middleware.RegisterJWTMiddleware()
	if err != nil {
		log.Fatalf("failed to initialize cognito jwt auth: %v", err)
	}

	adminAuth, err := middleware.RegisterAdminAuthMiddleware()
	if err != nil {
		log.Fatalf("failed to initialize cognito admin auth: %v", err)
	}

	// Register routes
	router.RegisterHealthCheckRoute(r)
	router.RegisterCategoryRoutes(r, db, jwtAuth, adminAuth)
	router.RegisterProductRoutes(r, db, jwtAuth, adminAuth)
	router.RegisterAttributeRoutes(r, db, jwtAuth, adminAuth)
	router.RegisterAttributeValueRoutes(r, db, jwtAuth, adminAuth)
	router.RegisterVariantRoutes(r, db, jwtAuth, adminAuth)
	router.RegisterProductImageRoutes(r, db, jwtAuth, adminAuth)

	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
