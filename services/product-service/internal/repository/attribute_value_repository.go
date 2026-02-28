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

type AttributeValueRepository struct {
	db *gorm.DB
}

func NewAttributeValueRepository(db *gorm.DB) *AttributeValueRepository {
	return &AttributeValueRepository{db: db}
}

func (r *AttributeValueRepository) CheckIfAttributeValueExists(id string) bool {
	var av entity.AttributeValue
	err := r.db.First(&av, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}
	return true
}

func (r *AttributeValueRepository) CheckIfAttributeExists(attributeID uint) bool {
	var a entity.Attribute
	err := r.db.First(&a, "id = ?", attributeID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}
	return true
}

func (r *AttributeValueRepository) GetAllAttributeValues() ([]entity.AttributeValue, error) {
	var avs []entity.AttributeValue
	if err := r.db.Preload("Attribute").Find(&avs).Error; err != nil {
		return nil, err
	}

	return avs, nil
}

func (r *AttributeValueRepository) GetAttributeValueByID(id string) (*response.AttributeValueResponse, error) {
	var av entity.AttributeValue
	if err := r.db.First(&av, "id = ?", id).Error; err != nil {
		return nil, err
	}

	avr := &response.AttributeValueResponse{
		ID:          av.ID,
		AttributeID: av.AttributeID,
		Value:       av.Value,
	}
	return avr, nil
}

func (r *AttributeValueRepository) AddAttributeValue(avReq *request.AttributeValueCreateRequest) error {
	if !r.CheckIfAttributeExists(avReq.AttributeID) {
		return pkg.AttributeNotFound
	}

	av := entity.AttributeValue{
		AttributeID: avReq.AttributeID,
		Value:       pkg.CapitalizeFirstLetter(avReq.Value),
	}
	if err := r.db.Create(&av).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return pkg.DuplicateEntry
		}
		return err
	}
	return nil
}

func (r *AttributeValueRepository) UpdateAttributeValue(id string, avReq *request.AttributeValueUpdateRequest) (*response.AttributeValueResponse, error) {
	if !r.CheckIfAttributeValueExists(id) {
		return nil, pkg.AttributeValueNotFound
	}

	updates := map[string]interface{}{}
	updates["value"] = pkg.CapitalizeFirstLetter(avReq.Value)

	if len(updates) == 0 {
		return nil, pkg.NoFieldsToUpdate
	}

	if err := r.db.Model(&entity.AttributeValue{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, pkg.DuplicateEntry
		}
		return nil, err
	}

	return r.GetAttributeValueByID(id)
}

func (r *AttributeValueRepository) DeleteAttributeValue(id string) error {
	if !r.CheckIfAttributeValueExists(id) {
		return pkg.AttributeValueNotFound
	}

	return r.db.Delete(&entity.AttributeValue{}, "id = ?", id).Error
}
