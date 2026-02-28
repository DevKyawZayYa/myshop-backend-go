package usecase

import (
	"context"
	"product-service/internal/entity"
	"product-service/internal/params"
	"product-service/internal/repository"
	"product-service/internal/request"
)

type ProductUsecase struct {
	productRepo *repository.ProductRepository
}

func NewProductUsecase(productRepo *repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context, params *params.ProductQueryParam) ([]*entity.Product, error) {
	return u.productRepo.GetAllProducts(ctx, params)
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	return u.productRepo.GetProductByID(ctx, id)
}

func (u *ProductUsecase) AddProduct(ctx context.Context, product *request.ProductCreateRequest) error {
	return u.productRepo.AddProduct(ctx, product)
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, id string, product *request.ProductUpdateRequest) (*entity.Product, error) {
	return u.productRepo.UpdateProduct(ctx, id, product)
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.productRepo.DeleteProduct(ctx, id)
}
