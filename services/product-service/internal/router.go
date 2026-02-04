package router

import (
	"net/http"
	"product-service/internal/handler"
	"product-service/internal/repository"
	"product-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHealthCheckRoute(r *gin.Engine) {
	r.GET("/product/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Product service is running"})
	})
}

func RegisterCategoryRoutes(r *gin.Engine, db *gorm.DB) {
	categoryGroup := r.Group("/product/categories")

	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)

	categoryGroup.GET("", categoryHandler.GetAllCategories)
	categoryGroup.GET("/:id", categoryHandler.GetCategoryByID)
	categoryGroup.POST("", categoryHandler.AddCategory)
	categoryGroup.PATCH("/:id", categoryHandler.PatchCategory)
	categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	categoryGroup.GET("/tree", categoryHandler.GetCategoryTree)
	categoryGroup.GET("/leaf", categoryHandler.GetLeafCategories)
	categoryGroup.GET("/:id/children", categoryHandler.GetChildCategoriesByID)
}

func RegisterProductRoutes(r *gin.Engine, db *gorm.DB) {
	productGroup := r.Group("/product/products")

	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	productGroup.GET("", productHandler.GetAllProducts)
	productGroup.GET("/:id", productHandler.GetProductByID)
	productGroup.POST("", productHandler.AddProduct)
	productGroup.PATCH("/:id", productHandler.PatchProduct)
	productGroup.DELETE("/:id", productHandler.DeleteProduct)
}

func RegisterAttributeRoutes(r *gin.Engine, db *gorm.DB) {
	attributeGroup := r.Group("/product/attributes")

	attributeRepo := repository.NewAttributeRepository(db)
	attributeUsecase := usecase.NewAttributeUsecase(attributeRepo)
	attributeHandler := handler.NewAttributeHandler(attributeUsecase)

	attributeGroup.GET("", attributeHandler.GetAllAttributes)
	attributeGroup.GET("/:id", attributeHandler.GetAttributeByID)
	attributeGroup.POST("", attributeHandler.AddAttribute)
	attributeGroup.PATCH("/:id", attributeHandler.PatchAttribute)
	attributeGroup.DELETE("/:id", attributeHandler.DeleteAttribute)
}

func RegisterAttributeValueRoutes(r *gin.Engine, db *gorm.DB) {
	attributeValueGroup := r.Group("/product/attribute-values")

	attributeValueRepo := repository.NewAttributeValueRepository(db)
	attributeValueUsecase := usecase.NewAttributeValueUsecase(attributeValueRepo)
	attributeValueHandler := handler.NewAttributeValueHandler(attributeValueUsecase)

	attributeValueGroup.GET("", attributeValueHandler.GetAllAttributeValues)
	attributeValueGroup.GET("/:id", attributeValueHandler.GetAttributeValueByID)
	attributeValueGroup.POST("", attributeValueHandler.AddAttributeValue)
	attributeValueGroup.PATCH("/:id", attributeValueHandler.PatchAttributeValue)
	attributeValueGroup.DELETE("/:id", attributeValueHandler.DeleteAttributeValue)
}

func RegisterVariantRoutes(r *gin.Engine, db *gorm.DB) {
	variantGroup := r.Group("/product/variants")

	variantRepo := repository.NewVariantRepository(db)
	variantUsecase := usecase.NewVariantUsecase(variantRepo)
	variantHandler := handler.NewVariantHandler(variantUsecase)

	variantGroup.GET("", variantHandler.GetAllVariants)
	variantGroup.GET("/:id", variantHandler.GetVariantByID)
	variantGroup.POST("", variantHandler.AddVariant)
	variantGroup.PATCH("/:id", variantHandler.PatchVariant)
	variantGroup.DELETE("/:id", variantHandler.DeleteVariant)
}

func RegisterProductImageRoutes(r *gin.Engine, db *gorm.DB) {
	productImageGroup := r.Group("/product/product-images")

	productImageRepo := repository.NewProductImageRepository(db)
	productImageUsecase := usecase.NewProductImageUsecase(productImageRepo)
	productImageHandler := handler.NewProductImageHandler(productImageUsecase)

	productImageGroup.GET("", productImageHandler.GetAllProductImages)
	productImageGroup.GET("/:id", productImageHandler.GetProductImageByID)
	productImageGroup.GET("/product/:productId", productImageHandler.GetImagesByProductID)
	productImageGroup.GET("/variant/:variantId", productImageHandler.GetImagesByVariantID)
	productImageGroup.POST("", productImageHandler.AddProductImage)
	productImageGroup.PATCH("/:id", productImageHandler.UpdateProductImage)
	productImageGroup.DELETE("/:id", productImageHandler.DeleteProductImage)
}
