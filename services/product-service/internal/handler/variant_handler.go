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

type VariantHandler struct {
	variantUsecase *usecase.VariantUsecase
}

func NewVariantHandler(variantUsecase *usecase.VariantUsecase) *VariantHandler {
	return &VariantHandler{variantUsecase: variantUsecase}
}

func (h *VariantHandler) GetAllVariants(c *gin.Context) {
	variantQueryParams := params.NewVariantQueryParam()
	if err := c.ShouldBindQuery(variantQueryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters"})
		return
	}

	variants, err := h.variantUsecase.GetAllVariants(c.Request.Context(), variantQueryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve variants"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": variants})
}

func (h *VariantHandler) GetVariantByID(c *gin.Context) {
	id := c.Param("id")
	variant, err := h.variantUsecase.GetVariantByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pkg.VariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve variant"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": variant})
}

func (h *VariantHandler) AddVariant(c *gin.Context) {
	var variant request.VariantRequest
	if err := c.ShouldBindJSON(&variant); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	if err := h.variantUsecase.AddVariant(c.Request.Context(), &variant); err != nil {
		switch {
		case errors.Is(err, pkg.ProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		case errors.Is(err, pkg.AttributeValueNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "One or more attribute values not found"})
		case errors.Is(err, pkg.InvalidAttributeValueForProduct):
			c.JSON(http.StatusBadRequest, gin.H{"message": "One or more attribute values do not belong to this product's attributes"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Variant with this SKU already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create variant"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Variant created successfully", "data": variant})
}

func (h *VariantHandler) PatchVariant(c *gin.Context) {
	id := c.Param("id")

	var variant request.VariantPatchRequest
	if err := c.ShouldBindJSON(&variant); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	updatedVariant, err := h.variantUsecase.UpdateVariant(c.Request.Context(), id, &variant)
	if err != nil {
		switch {
		case errors.Is(err, pkg.VariantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Variant not found"})
		case errors.Is(err, pkg.ProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		case errors.Is(err, pkg.AttributeValueNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "One or more attribute values not found"})
		case errors.Is(err, pkg.InvalidAttributeValueForProduct):
			c.JSON(http.StatusBadRequest, gin.H{"message": "One or more attribute values do not belong to this product's attributes"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Variant with this SKU already exists"})
		case errors.Is(err, pkg.NoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields provided to update"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update variant"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variant updated successfully", "data": updatedVariant})
}

func (h *VariantHandler) DeleteVariant(c *gin.Context) {
	id := c.Param("id")
	if err := h.variantUsecase.DeleteVariant(c.Request.Context(), id); err != nil {
		if errors.Is(err, pkg.VariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete variant"})
		return
	}
	c.Status(http.StatusNoContent)
}
