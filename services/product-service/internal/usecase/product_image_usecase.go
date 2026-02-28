package usecase

import (
	"context"
	"product-service/internal/entity"
	"product-service/internal/repository"
	"product-service/internal/request"
)

type ProductImageUsecase struct {
	productImageRepo *repository.ProductImageRepository
}

func NewProductImageUsecase(productImageRepo *repository.ProductImageRepository) *ProductImageUsecase {
	return &ProductImageUsecase{productImageRepo: productImageRepo}
}

func (uc *ProductImageUsecase) GetAllProductImages(ctx context.Context) ([]*entity.ProductImage, error) {
	return uc.productImageRepo.GetAllProductImages(ctx)
}

func (uc *ProductImageUsecase) GetProductImageByID(ctx context.Context, id string) (*entity.ProductImage, error) {
	return uc.productImageRepo.GetProductImageByID(ctx, id)
}

func (uc *ProductImageUsecase) GetImagesByProductID(ctx context.Context, productID uint) ([]*entity.ProductImage, error) {
	return uc.productImageRepo.GetImagesByProductID(ctx, productID)
}

func (uc *ProductImageUsecase) GetImagesByVariantID(ctx context.Context, variantID uint) ([]*entity.ProductImage, error) {
	return uc.productImageRepo.GetImagesByVariantID(ctx, variantID)
}

func (uc *ProductImageUsecase) AddProductImage(ctx context.Context, imageReq *request.ProductImageCreateRequest) (*entity.ProductImage, error) {
	return uc.productImageRepo.AddProductImage(ctx, imageReq)
}

func (uc *ProductImageUsecase) UpdateProductImage(ctx context.Context, id string, imageReq *request.ProductImageUpdateRequest) (*entity.ProductImage, error) {
	return uc.productImageRepo.UpdateProductImage(ctx, id, imageReq)
}

func (uc *ProductImageUsecase) DeleteProductImage(ctx context.Context, id string) error {
	return uc.productImageRepo.DeleteProductImage(ctx, id)
}

func (uc *ProductImageUsecase) AddProductImageBatch(ctx context.Context, images []*request.ProductImageCreateRequest) ([]*entity.ProductImage, error) {
	return uc.productImageRepo.AddProductImageBatch(ctx, images)
}
