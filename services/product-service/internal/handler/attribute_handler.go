package handler

import (
	"errors"
	"net/http"
	"product-service/internal/request"
	"product-service/internal/usecase"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

type AttributeHandler struct {
	attributeUsecase *usecase.AttributeUsecase
}

func NewAttributeHandler(attributeUsecase *usecase.AttributeUsecase) *AttributeHandler {
	return &AttributeHandler{attributeUsecase: attributeUsecase}
}

func (h *AttributeHandler) GetAllAttributes(c *gin.Context) {
	attrs, err := h.attributeUsecase.GetAllAttributes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve attributes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": attrs})
}

func (h *AttributeHandler) AddAttribute(c *gin.Context) {
	var attr request.AttributeRequest
	if err := c.ShouldBindJSON(&attr); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	if err := h.attributeUsecase.AddAttribute(&attr); err != nil {
		switch {
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Attribute name already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create attribute"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Attribute created successfully", "data": attr})
}

func (h *AttributeHandler) GetAttributeByID(c *gin.Context) {
	id := c.Param("id")
	attr, err := h.attributeUsecase.GetAttributeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attr)
}

func (h *AttributeHandler) PatchAttribute(c *gin.Context) {
	id := c.Param("id")

	var attr request.AttributePatchRequest
	if err := c.ShouldBindJSON(&attr); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	updated, err := h.attributeUsecase.UpdateAttribute(id, &attr)
	if err != nil {
		switch {
		case errors.Is(err, pkg.AttributeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Attribute name already exists"})
		case errors.Is(err, pkg.NoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields provided to update"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update attribute"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attribute updated successfully", "data": updated})
}

func (h *AttributeHandler) DeleteAttribute(c *gin.Context) {
	id := c.Param("id")
	if err := h.attributeUsecase.DeleteAttribute(id); err != nil {
		switch {
		case errors.Is(err, pkg.AttributeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found"})
		case errors.Is(err, pkg.AttributeHasValues):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Attribute has values and cannot be deleted"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
