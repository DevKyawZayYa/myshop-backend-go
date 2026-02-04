package repository

import (
	"errors"
	"product-service/internal/entity"
	"product-service/internal/request"
	"product-service/internal/response"

	"myshop-shared/pkg"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories() ([]response.CategoryResponse, error) {
	var categories []response.CategoryResponse
	if err := r.db.Model(&entity.Category{}).
		Select("id, name, description, parent_id").
		Scan(&categories).Error; err != nil {
		return nil, err
	}
	if categories == nil {
		categories = make([]response.CategoryResponse, 0)
	}
	return categories, nil
}

func (r *CategoryRepository) CheckIfCategoryExists(id string) bool {
	var category entity.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		// For other errors, log or handle as needed. Here, treat as not found.
		return false
	}
	return true
}

func (r *CategoryRepository) GetCategoryByID(id string) (*response.CategoryResponse, error) {
	var category entity.Category
	if err := r.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
	}, nil
}

func (r *CategoryRepository) GetChildCategoriesByID(id string) ([]response.CategoryResponse, error) {
	var categories []response.CategoryResponse
	if err := r.db.Model(&entity.Category{}).
		Select("id, name, description, parent_id").
		Where("parent_id = ?", id).
		Scan(&categories).Error; err != nil {
		return nil, err
	}
	if categories == nil {
		categories = make([]response.CategoryResponse, 0)
	}
	return categories, nil
}

func (r *CategoryRepository) AddCategory(category *request.CategoryRequest) error {
	if category.ParentID != nil {
		if !r.CheckIfCategoryExists(pkg.UintToString(*category.ParentID)) {
			return pkg.ParentCategoryNotFound
		}
	}
	cat := entity.Category{
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
	}
	err := r.db.Create(&cat).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return pkg.DuplicateEntry
		}
		return err
	}
	return nil
}

func (r *CategoryRepository) UpdateCategory(id string, category *request.CategoryPatchRequest) (
	*response.CategoryResponse, error) {
	if !r.CheckIfCategoryExists(id) {
		return nil, pkg.CategoryNotFound
	}

	if category.ParentID != nil {
		if *category.ParentID == pkg.StringToUint(id) {
			return nil, pkg.CategoryCannotBeItsOwnParent
		}
		if !r.CheckIfCategoryExists(pkg.UintToString(*category.ParentID)) {
			return nil, pkg.ParentCategoryNotFound
		}
	}

	if err := r.db.Model(&entity.Category{}).Where("id = ?", id).Updates(category).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, pkg.DuplicateEntry
		}
		return nil, err
	}

	return nil, nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	if !r.CheckIfCategoryExists(id) {
		return pkg.CategoryNotFound
	}
	// Check if category is parent of any child category
	var childCount int64
	if err := r.db.Model(&entity.Category{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}

	if childCount > 0 {
		return pkg.CategoryHasChildren // You should define this error in your pkg
	}

	return r.db.Delete(&entity.Category{}, "id = ?", id).Error
}

func (r *CategoryRepository) GetCategoryTree() ([]*response.CategoryTreeResponse, error) {
	// Fetch all categories in one query
	var allCategories []response.CategoryResponse
	if err := r.db.Model(&entity.Category{}).
		Select("id, name, description, parent_id").
		Scan(&allCategories).Error; err != nil {
		return nil, err
	}

	// Build maps for quick lookup
	categoryMap := make(map[uint]*response.CategoryTreeResponse)
	var roots []*response.CategoryTreeResponse

	// First pass: create all tree nodes
	for i := range allCategories {
		categoryMap[allCategories[i].ID] = &response.CategoryTreeResponse{
			CategoryResponse: allCategories[i],
			Children:         []*response.CategoryTreeResponse{},
		}
	}

	// Second pass: build parent-child relationships
	for i := range allCategories {
		node := categoryMap[allCategories[i].ID]
		if allCategories[i].ParentID == nil {
			// Root category
			roots = append(roots, node)
		} else {
			// Add to parent's children
			if parent, exists := categoryMap[*allCategories[i].ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots, nil
}

func (r *CategoryRepository) GetLeafCategories() ([]response.CategoryResponse, error) {
	var categories []response.CategoryResponse
	if err := r.db.Model(&entity.Category{}).
		Select("id, name, description, parent_id").
		Where("id NOT IN (?)",
			r.db.Model(&entity.Category{}).
				Select("DISTINCT parent_id").
				Where("parent_id IS NOT NULL"),
		).
		Scan(&categories).Error; err != nil {
		return nil, err
	}
	if categories == nil {
		categories = make([]response.CategoryResponse, 0)
	}
	return categories, nil
}
