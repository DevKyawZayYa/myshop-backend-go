package handler

import (
	"errors"
	"net/http"
	"product-service/internal/request"
	"product-service/internal/usecase"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

type AttributeValueHandler struct {
	attributeValueUsecase *usecase.AttributeValueUsecase
}

func NewAttributeValueHandler(attributeValueUsecase *usecase.AttributeValueUsecase) *AttributeValueHandler {
	return &AttributeValueHandler{attributeValueUsecase: attributeValueUsecase}
}

func (h *AttributeValueHandler) GetAllAttributeValues(c *gin.Context) {
	avs, err := h.attributeValueUsecase.GetAllAttributeValues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve attribute values"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": avs})
}

func (h *AttributeValueHandler) AddAttributeValue(c *gin.Context) {
	var avReq request.AttributeValueRequest
	if err := c.ShouldBindJSON(&avReq); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	if err := h.attributeValueUsecase.AddAttributeValue(&avReq); err != nil {
		switch {
		case errors.Is(err, pkg.AttributeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Attribute value already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create attribute value"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Attribute value created successfully", "data": avReq})
}

func (h *AttributeValueHandler) GetAttributeValueByID(c *gin.Context) {
	id := c.Param("id")
	av, err := h.attributeValueUsecase.GetAttributeValueByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, av)
}

func (h *AttributeValueHandler) PatchAttributeValue(c *gin.Context) {
	id := c.Param("id")

	var avReq request.AttributeValuePatchRequest
	if err := c.ShouldBindJSON(&avReq); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	updated, err := h.attributeValueUsecase.UpdateAttributeValue(id, &avReq)
	if err != nil {
		switch {
		case errors.Is(err, pkg.AttributeValueNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute value not found"})
		case errors.Is(err, pkg.AttributeNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute not found"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Attribute value already exists"})
		case errors.Is(err, pkg.NoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields provided to update"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update attribute value"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attribute value updated successfully", "data": updated})
}

func (h *AttributeValueHandler) DeleteAttributeValue(c *gin.Context) {
	id := c.Param("id")
	if err := h.attributeValueUsecase.DeleteAttributeValue(id); err != nil {
		switch {
		case errors.Is(err, pkg.AttributeValueNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Attribute value not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
