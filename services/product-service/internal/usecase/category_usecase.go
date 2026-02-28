package usecase

import (
	"product-service/internal/repository"
	"product-service/internal/request"
	"product-service/internal/response"
)

type CategoryUsecase struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo *repository.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{categoryRepo: categoryRepo}
}

func (u *CategoryUsecase) GetAllCategories() ([]response.CategoryResponse, error) {
	return u.categoryRepo.GetAllCategories()
}

func (u *CategoryUsecase) AddCategory(category *request.CategoryCreateRequest) error {
	return u.categoryRepo.AddCategory(category)
}

func (u *CategoryUsecase) UpdateCategory(id string, category *request.CategoryUpdateRequest) (*response.CategoryResponse, error) {
	return u.categoryRepo.UpdateCategory(id, category)
}

func (u *CategoryUsecase) GetCategoryByID(id string) (*response.CategoryResponse, error) {
	return u.categoryRepo.GetCategoryByID(id)
}

func (u *CategoryUsecase) DeleteCategory(id string) error {
	return u.categoryRepo.DeleteCategory(id)
}

func (u *CategoryUsecase) GetCategoryTree() ([]*response.CategoryTreeResponse, error) {
	return u.categoryRepo.GetCategoryTree()
}

func (u *CategoryUsecase) GetChildCategoriesByID(id string) ([]response.CategoryResponse, error) {
	return u.categoryRepo.GetChildCategoriesByID(id)
}

func (u *CategoryUsecase) GetLeafCategories() ([]response.CategoryResponse, error) {
	return u.categoryRepo.GetLeafCategories()
}
