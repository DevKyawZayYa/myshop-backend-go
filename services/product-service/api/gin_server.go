package server

import (
	"fmt"
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

	// Register routes
	router.RegisterHealthCheckRoute(r)
	router.RegisterCategoryRoutes(r, db)
	router.RegisterProductRoutes(r, db)
	router.RegisterAttributeRoutes(r, db)
	router.RegisterAttributeValueRoutes(r, db)
	router.RegisterVariantRoutes(r, db)
	router.RegisterProductImageRoutes(r, db)

	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
