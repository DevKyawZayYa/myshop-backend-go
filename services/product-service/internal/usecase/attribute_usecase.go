package usecase

import (
	"product-service/internal/entity"
	"product-service/internal/repository"
	"product-service/internal/request"
)

type AttributeUsecase struct {
	attributeRepo *repository.AttributeRepository
}

func NewAttributeUsecase(attributeRepo *repository.AttributeRepository) *AttributeUsecase {
	return &AttributeUsecase{attributeRepo: attributeRepo}
}

func (u *AttributeUsecase) GetAllAttributes() ([]entity.Attribute, error) {
	return u.attributeRepo.GetAllAttributes()
}

func (u *AttributeUsecase) AddAttribute(attr *request.AttributeCreateRequest) error {
	return u.attributeRepo.AddAttribute(attr)
}

func (u *AttributeUsecase) GetAttributeByID(id string) (*entity.Attribute, error) {
	return u.attributeRepo.GetAttributeByID(id)
}

func (u *AttributeUsecase) UpdateAttribute(id string, attr *request.AttributeUpdateRequest) (*entity.Attribute, error) {
	return u.attributeRepo.UpdateAttribute(id, attr)
}

func (u *AttributeUsecase) DeleteAttribute(id string) error {
	return u.attributeRepo.DeleteAttribute(id)
}
