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

func RegisterCategoryRoutes(r *gin.Engine, db *gorm.DB, jwtAuth gin.HandlerFunc, adminAuth gin.HandlerFunc) {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)

	// Public category routes (requires JWT)
	publicCategoryGroup := r.Group("/product/categories", jwtAuth)
	{
		// Get hierarchical category tree for frontend navigation
		publicCategoryGroup.GET("/tree", categoryHandler.GetCategoryTree)

		// Get leaf categories (categories without children)
		publicCategoryGroup.GET("/leaf", categoryHandler.GetLeafCategories)

		// Get specific category by ID
		publicCategoryGroup.GET("/:id", categoryHandler.GetCategoryByID)

		// Get child categories of a specific category
		publicCategoryGroup.GET("/:id/children", categoryHandler.GetChildCategoriesByID)
	}

	// Admin category routes (requires admin token)
	adminCategoryGroup := r.Group("/product/admin/categories", adminAuth)
	{
		// Create new category
		adminCategoryGroup.POST("", categoryHandler.AddCategory)

		// Update category
		adminCategoryGroup.PATCH("/:id", categoryHandler.PatchCategory)

		// Delete category
		adminCategoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}

func RegisterProductRoutes(r *gin.Engine, db *gorm.DB, jwtAuth gin.HandlerFunc, adminAuth gin.HandlerFunc) {
	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase)

	// Public product routes (requires JWT)
	publicProductGroup := r.Group("/product/products", jwtAuth)
	{
		// Get all products with filtering and pagination
		publicProductGroup.GET("", productHandler.GetAllProducts)

		// Get specific product by ID
		publicProductGroup.GET("/:id", productHandler.GetProductByID)
	}

	// Admin product routes (requires admin token)
	adminProductGroup := r.Group("/product/admin/products", adminAuth)
	{
		// Create new product
		adminProductGroup.POST("", productHandler.AddProduct)

		// Update product
		adminProductGroup.PATCH("/:id", productHandler.PatchProduct)

		// Delete product
		adminProductGroup.DELETE("/:id", productHandler.DeleteProduct)
	}
}

func RegisterAttributeRoutes(r *gin.Engine, db *gorm.DB, jwtAuth gin.HandlerFunc, adminAuth gin.HandlerFunc) {
	attributeRepo := repository.NewAttributeRepository(db)
	attributeUsecase := usecase.NewAttributeUsecase(attributeRepo)
	attributeHandler := handler.NewAttributeHandler(attributeUsecase)

	// Public attribute routes (requires JWT)
	publicAttributeGroup := r.Group("/product/attributes", jwtAuth)
	{
		// Get all attributes
		publicAttributeGroup.GET("", attributeHandler.GetAllAttributes)

		// Get specific attribute by ID
		publicAttributeGroup.GET("/:id", attributeHandler.GetAttributeByID)
	}

	// Admin attribute routes (requires admin token)
	adminAttributeGroup := r.Group("/product/admin/attributes", adminAuth)
	{
		// Create new attribute
		adminAttributeGroup.POST("", attributeHandler.AddAttribute)

		// Update attribute
		adminAttributeGroup.PATCH("/:id", attributeHandler.PatchAttribute)

		// Delete attribute
		adminAttributeGroup.DELETE("/:id", attributeHandler.DeleteAttribute)
	}
}

func RegisterAttributeValueRoutes(r *gin.Engine, db *gorm.DB, jwtAuth gin.HandlerFunc, adminAuth gin.HandlerFunc) {
	attributeValueRepo := repository.NewAttributeValueRepository(db)
	attributeValueUsecase := usecase.NewAttributeValueUsecase(attributeValueRepo)
	attributeValueHandler := handler.NewAttributeValueHandler(attributeValueUsecase)

	// Public attribute value routes (requires JWT)
	publicAttributeValueGroup := r.Group("/product/attribute-values", jwtAuth)
	{
		// Get all attribute values
		publicAttributeValueGroup.GET("", attributeValueHandler.GetAllAttributeValues)

		// Get specific attribute value by ID
		publicAttributeValueGroup.GET("/:id", attributeValueHandler.GetAttributeValueByID)
	}

	// Admin attribute value routes (requires admin token)
	adminAttributeValueGroup := r.Group("/product/admin/attribute-values", adminAuth)
	{
		// Create new attribute value
		adminAttributeValueGroup.POST("", attributeValueHandler.AddAttributeValue)

		// Update attribute value
		adminAttributeValueGroup.PATCH("/:id", attributeValueHandler.PatchAttributeValue)

		// Delete attribute value
		adminAttributeValueGroup.DELETE("/:id", attributeValueHandler.DeleteAttributeValue)
	}
}

func RegisterVariantRoutes(r *gin.Engine, db *gorm.DB, jwtAuth gin.HandlerFunc, adminAuth gin.HandlerFunc) {
	variantRepo := repository.NewVariantRepository(db)
	variantUsecase := usecase.NewVariantUsecase(variantRepo)
	variantHandler := handler.NewVariantHandler(variantUsecase)

	// Public variant routes (requires JWT)
	publicVariantGroup := r.Group("/product/variants", jwtAuth)
	{
		// Get all variants with filtering and pagination
		publicVariantGroup.GET("", variantHandler.GetAllVariants)

		// Get specific variant by ID
		publicVariantGroup.GET("/:id", variantHandler.GetVariantByID)
	}

	// Admin variant routes (requires admin token)
	adminVariantGroup := r.Group("/product/admin/variants", adminAuth)
	{
		// Create new variant
		adminVariantGroup.POST("", variantHandler.AddVariant)

		// Update variant
		adminVariantGroup.PATCH("/:id", variantHandler.PatchVariant)

		// Delete variant
		adminVariantGroup.DELETE("/:id", variantHandler.DeleteVariant)
	}
}

func RegisterProductImageRoutes(r *gin.Engine, db *gorm.DB, jwtAuth gin.HandlerFunc, adminAuth gin.HandlerFunc) {
	productImageRepo := repository.NewProductImageRepository(db)
	productImageUsecase := usecase.NewProductImageUsecase(productImageRepo)
	productImageHandler := handler.NewProductImageHandler(productImageUsecase)

	// Public product image routes (requires JWT)
	publicImageGroup := r.Group("/product/product-images", jwtAuth)
	{
		// Get all product images
		publicImageGroup.GET("", productImageHandler.GetAllProductImages)

		// Get specific product image by ID
		publicImageGroup.GET("/:id", productImageHandler.GetProductImageByID)

		// Get images for a specific product
		publicImageGroup.GET("/product/:productId", productImageHandler.GetImagesByProductID)

		// Get images for a specific variant
		publicImageGroup.GET("/variant/:variantId", productImageHandler.GetImagesByVariantID)
	}

	// Admin product image routes (requires admin token)
	adminImageGroup := r.Group("/product/admin/product-images", adminAuth)
	{
		// Add new product image
		adminImageGroup.POST("", productImageHandler.AddProductImage)

		// Update product image
		adminImageGroup.PATCH("/:id", productImageHandler.UpdateProductImage)

		// Delete product image
		adminImageGroup.DELETE("/:id", productImageHandler.DeleteProductImage)
	}
}
