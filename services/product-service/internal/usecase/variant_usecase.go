package usecase

import (
	"context"
	"product-service/internal/entity"
	"product-service/internal/params"
	"product-service/internal/repository"
	"product-service/internal/request"
)

type VariantUsecase struct {
	variantRepo *repository.VariantRepository
}

func NewVariantUsecase(variantRepo *repository.VariantRepository) *VariantUsecase {
	return &VariantUsecase{variantRepo: variantRepo}
}

func (u *VariantUsecase) GetAllVariants(ctx context.Context, params *params.VariantQueryParam) ([]*entity.Variant, error) {
	return u.variantRepo.GetAllVariants(ctx, params)
}

func (u *VariantUsecase) GetVariantByID(ctx context.Context, id string) (*entity.Variant, error) {
	return u.variantRepo.GetVariantByID(ctx, id)
}

func (u *VariantUsecase) AddVariant(ctx context.Context, variant *request.VariantCreateRequest) error {
	return u.variantRepo.AddVariant(ctx, variant)
}

func (u *VariantUsecase) UpdateVariant(ctx context.Context, id string, variant *request.VariantUpdateRequest) (*entity.Variant, error) {
	return u.variantRepo.UpdateVariant(ctx, id, variant)
}

func (u *VariantUsecase) DeleteVariant(ctx context.Context, id string) error {
	return u.variantRepo.DeleteVariant(ctx, id)
}
