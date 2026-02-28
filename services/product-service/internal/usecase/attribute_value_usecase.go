package usecase

import (
	"product-service/internal/entity"
	"product-service/internal/repository"
	"product-service/internal/request"
	"product-service/internal/response"
)

type AttributeValueUsecase struct {
	attributeValueRepo *repository.AttributeValueRepository
}

func NewAttributeValueUsecase(attributeValueRepo *repository.AttributeValueRepository) *AttributeValueUsecase {
	return &AttributeValueUsecase{attributeValueRepo: attributeValueRepo}
}

func (u *AttributeValueUsecase) GetAllAttributeValues() ([]entity.AttributeValue, error) {
	return u.attributeValueRepo.GetAllAttributeValues()
}

func (u *AttributeValueUsecase) AddAttributeValue(avReq *request.AttributeValueCreateRequest) error {
	return u.attributeValueRepo.AddAttributeValue(avReq)
}

func (u *AttributeValueUsecase) GetAttributeValueByID(id string) (*response.AttributeValueResponse, error) {
	return u.attributeValueRepo.GetAttributeValueByID(id)
}

func (u *AttributeValueUsecase) UpdateAttributeValue(id string, avReq *request.AttributeValueUpdateRequest) (*response.AttributeValueResponse, error) {
	return u.attributeValueRepo.UpdateAttributeValue(id, avReq)
}

func (u *AttributeValueUsecase) DeleteAttributeValue(id string) error {
	return u.attributeValueRepo.DeleteAttributeValue(id)
}
