package handler

import (
	"errors"
	"net/http"
	"product-service/internal/request"
	"product-service/internal/usecase"

	"myshop-shared/pkg"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUsecase *usecase.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase}
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryUsecase.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CategoryHandler) AddCategory(c *gin.Context) {
	var category request.CategoryRequest
	if err := c.ShouldBindJSON(&category); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	if err := h.categoryUsecase.AddCategory(&category); err != nil {
		switch {
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Category name already exists"})
		case errors.Is(err, pkg.ParentCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Parent category not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create category"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "data": category})
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryUsecase.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) PatchCategory(c *gin.Context) {
	id := c.Param("id")

	var category request.CategoryPatchRequest
	if err := c.ShouldBindJSON(&category); err != nil {
		pkg.HandleValidationError(c, err)
		return
	}

	updatedCategory, err := h.categoryUsecase.UpdateCategory(id, &category)
	if err != nil {
		switch {
		case errors.Is(err, pkg.CategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		case errors.Is(err, pkg.ParentCategoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Parent category not found"})
		case errors.Is(err, pkg.CategoryCannotBeItsOwnParent):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Category cannot be its own parent"})
		case errors.Is(err, pkg.DuplicateEntry):
			c.JSON(http.StatusConflict, gin.H{"message": "Category name already exists"})
		case errors.Is(err, pkg.NoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"message": "No fields provided to update"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update category"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "data": updatedCategory})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.categoryUsecase.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := h.categoryUsecase.GetCategoryTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve category tree"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CategoryHandler) GetChildCategoriesByID(c *gin.Context) {
	categories, err := h.categoryUsecase.GetChildCategoriesByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve child categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CategoryHandler) GetLeafCategories(c *gin.Context) {
	categories, err := h.categoryUsecase.GetLeafCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve leaf categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}
