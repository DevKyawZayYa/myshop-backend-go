package repository

import (
	"errors"
	"product-service/internal/entity"
	"product-service/internal/request"

	"myshop-shared/pkg"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type AttributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) *AttributeRepository {
	return &AttributeRepository{db: db}
}

func (r *AttributeRepository) CheckIfAttributeExists(id string) bool {
	var a entity.Attribute
	err := r.db.First(&a, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}
	return true
}

func (r *AttributeRepository) GetAllAttributes() ([]entity.Attribute, error) {
	var attrs []entity.Attribute
	if err := r.db.Preload("Values").Find(&attrs).Error; err != nil {
		return nil, err
	}

	return attrs, nil
}

func (r *AttributeRepository) GetAttributeByID(id string) (*entity.Attribute, error) {
	var a entity.Attribute
	if err := r.db.Preload("Values").First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *AttributeRepository) AddAttribute(attr *request.AttributeRequest) error {
	a := entity.Attribute{Name: pkg.CapitalizeFirstLetter(attr.Name)}
	if err := r.db.Create(&a).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return pkg.DuplicateEntry
		}
		return err
	}
	return nil
}

func (r *AttributeRepository) UpdateAttribute(id string, attr *request.AttributePatchRequest) (*entity.Attribute, error) {
	if !r.CheckIfAttributeExists(id) {
		return nil, pkg.AttributeNotFound
	}

	updates := map[string]interface{}{}
	if attr.Name != nil {
		updates["name"] = *attr.Name
	}

	if len(updates) == 0 {
		return nil, pkg.NoFieldsToUpdate
	}

	if err := r.db.Model(&entity.Attribute{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, pkg.DuplicateEntry
		}
		return nil, err
	}

	return r.GetAttributeByID(id)
}

func (r *AttributeRepository) DeleteAttribute(id string) error {
	if !r.CheckIfAttributeExists(id) {
		return pkg.AttributeNotFound
	}

	var cnt int64
	if err := r.db.Model(&entity.AttributeValue{}).Where("attribute_id = ?", id).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return pkg.AttributeHasValues
	}

	return r.db.Delete(&entity.Attribute{}, "id = ?", id).Error
}
