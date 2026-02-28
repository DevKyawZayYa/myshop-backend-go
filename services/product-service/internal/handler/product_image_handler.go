package handler

import (
	"errors"
	"net/http"
	"product-service/internal/request"
	"product-service/internal/response"
	"product-service/internal/usecase"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

type ProductImageHandler struct {
	productImageUsecase *usecase.ProductImageUsecase
}

func NewProductImageHandler(productImageUsecase *usecase.ProductImageUsecase) *ProductImageHandler {
	return &ProductImageHandler{productImageUsecase: productImageUsecase}
}

func (h *ProductImageHandler) GetAllProductImages(c *gin.Context) {
	images, err := h.productImageUsecase.GetAllProductImages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve product images"})
		return
	}

	responses := make([]response.ProductImageResponse, len(images))
	for i, img := range images {
		responses[i] = response.ProductImageResponse{
			ID:        img.ID,
			ProductID: img.ProductID,
			VariantID: img.VariantID,
			URL:       img.URL,
			IsDefault: img.IsDefault,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *ProductImageHandler) GetProductImageByID(c *gin.Context) {
	id := c.Param("id")
	image, err := h.productImageUsecase.GetProductImageByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pkg.ProductImageNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve product image"})
		return
	}

	resp := response.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		VariantID: image.VariantID,
		URL:       image.URL,
		IsDefault: image.IsDefault,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *ProductImageHandler) GetImagesByProductID(c *gin.Context) {
	productID := c.Param("productId")

	id := pkg.StringToUint(productID)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	images, err := h.productImageUsecase.GetImagesByProductID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve product images"})
		return
	}

	responses := make([]response.ProductImageResponse, len(images))
	for i, img := range images {
		responses[i] = response.ProductImageResponse{
			ID:        img.ID,
			ProductID: img.ProductID,
			VariantID: img.VariantID,
			URL:       img.URL,
			IsDefault: img.IsDefault,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *ProductImageHandler) GetImagesByVariantID(c *gin.Context) {
	variantID := c.Param("variantId")

	id := pkg.StringToUint(variantID)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid variant ID"})
		return
	}

	images, err := h.productImageUsecase.GetImagesByVariantID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve variant images"})
		return
	}

	responses := make([]response.ProductImageResponse, len(images))
	for i, img := range images {
		responses[i] = response.ProductImageResponse{
			ID:        img.ID,
			ProductID: img.ProductID,
			VariantID: img.VariantID,
			URL:       img.URL,
			IsDefault: img.IsDefault,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *ProductImageHandler) AddProductImage(c *gin.Context) {
	var imageReq request.ProductImageCreateRequest
	if err := c.ShouldBindJSON(&imageReq); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	image, err := h.productImageUsecase.AddProductImage(c.Request.Context(), &imageReq)
	if err != nil {
		if errors.Is(err, pkg.ProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}
		if errors.Is(err, pkg.VariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add product image"})
		return
	}

	resp := response.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		VariantID: image.VariantID,
		URL:       image.URL,
		IsDefault: image.IsDefault,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *ProductImageHandler) UpdateProductImage(c *gin.Context) {
	id := c.Param("id")

	var imageReq request.ProductImageUpdateRequest
	if err := c.ShouldBindJSON(&imageReq); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	image, err := h.productImageUsecase.UpdateProductImage(c.Request.Context(), id, &imageReq)
	if err != nil {
		if errors.Is(err, pkg.NoFieldsToUpdate) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields to update"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product image"})
		return
	}

	resp := response.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		VariantID: image.VariantID,
		URL:       image.URL,
		IsDefault: image.IsDefault,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *ProductImageHandler) DeleteProductImage(c *gin.Context) {
	id := c.Param("id")

	err := h.productImageUsecase.DeleteProductImage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product image deleted successfully"})
}
