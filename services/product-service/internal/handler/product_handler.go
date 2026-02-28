package handler

import (
	"errors"
	"net/http"
	"product-service/internal/params"
	"product-service/internal/request"
	"product-service/internal/usecase"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	productQueryParams := params.NewProductQueryParam()
	if err := c.ShouldBindQuery(productQueryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}

	products, err := h.productUsecase.GetAllProducts(c.Request.Context(), productQueryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productUsecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pkg.ProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var product request.ProductCreateRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	if err := h.productUsecase.AddProduct(c.Request.Context(), &product); err != nil {
		switch {
		case errors.Is(err, pkg.CategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "One or more categories not found"})
		case errors.Is(err, pkg.AttributeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "One or more attributes not found"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Product already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "data": product})
}

func (h *ProductHandler) PatchProduct(c *gin.Context) {
	id := c.Param("id")

	var product request.ProductUpdateRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	updatedProduct, err := h.productUsecase.UpdateProduct(c.Request.Context(), id, &product)
	if err != nil {
		switch {
		case errors.Is(err, pkg.ProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		case errors.Is(err, pkg.CategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "One or more categories not found"})
		case errors.Is(err, pkg.AttributeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "One or more attributes not found"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Product name already exists"})
		case errors.Is(err, pkg.NoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields provided to update"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "data": updatedProduct})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productUsecase.DeleteProduct(c.Request.Context(), id); err != nil {
		if errors.Is(err, pkg.ProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product"})
		return
	}
	c.Status(http.StatusNoContent)
}
