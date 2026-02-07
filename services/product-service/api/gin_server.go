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

	adminAuth, err := middleware.RegisterAuthMiddleware()

	if err != nil {
		log.Fatalf("failed to initialize cognito auth: %v", err)
	}

	// Register routes
	router.RegisterHealthCheckRoute(r)
	router.RegisterCategoryRoutes(r, db, adminAuth)
	router.RegisterProductRoutes(r, db, adminAuth)
	router.RegisterAttributeRoutes(r, db, adminAuth)
	router.RegisterAttributeValueRoutes(r, db, adminAuth)
	router.RegisterVariantRoutes(r, db, adminAuth)
	router.RegisterProductImageRoutes(r, db, adminAuth)

	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
