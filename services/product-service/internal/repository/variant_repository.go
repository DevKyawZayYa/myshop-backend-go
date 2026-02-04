package repository

import (
	"context"
	"errors"
	"product-service/internal/entity"
	"product-service/internal/params"
	"product-service/internal/request"

	"myshop-shared/pkg"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type VariantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

func (r *VariantRepository) GetAllVariants(ctx context.Context, params *params.VariantQueryParam) ([]*entity.Variant, error) {
	var variants []*entity.Variant

	query := r.db.WithContext(ctx).Model(&entity.Variant{})

	if params.ProductID != "" {
		query = query.Where("product_id = ?", params.ProductID)
	}

	query = pkg.NewQueryBuilder(query).
		ApplyPagination(params.Limit, params.Offset).
		ApplyPreload([]string{"Product", "ProductImages", "AttributeValues"}).
		Build()

	if err := query.Find(&variants).Error; err != nil {
		return nil, err
	}

	return variants, nil
}

func (r *VariantRepository) GetVariantByID(ctx context.Context, id string) (*entity.Variant, error) {
	var variant = new(entity.Variant)
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("ProductImages").
		Preload("AttributeValues").
		First(&variant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.VariantNotFound
		}
		return nil, err
	}

	return variant, nil
}

func (r *VariantRepository) AddVariant(ctx context.Context, variantReq *request.VariantRequest) error {
	// Check if product exists
	var product entity.Product
	if err := r.db.WithContext(ctx).Preload("Attributes").
		First(&product, "id = ?", variantReq.ProductID).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ProductNotFound
		}
		return err
	}

	// Check if attributes exist and belong to product
	var attributeValues []entity.AttributeValue
	if len(variantReq.AttributeValues) > 0 {
		attributeValues = r.CheckAttributeValuesExist(ctx, variantReq.AttributeValues)
		if attributeValues == nil {
			return pkg.AttributeValueNotFound
		}

		// Validate that all attribute values belong to product's attributeValues
		if !r.ValidateAttributeValuesForProduct(ctx, product, attributeValues) {
			return pkg.InvalidAttributeValueForProduct
		}
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		variant := entity.Variant{
			ProductID:       variantReq.ProductID,
			SKU:             variantReq.SKU,
			BasePrice:       variantReq.BasePrice,
			ComparePrice:    variantReq.ComparePrice,
			Stock:           variantReq.Stock,
			AttributeValues: attributeValues,
		}

		if err := tx.Create(&variant).Error; err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return pkg.DuplicateEntry
			}
			return err
		}

		// Add variant images if provided
		if len(variantReq.ProductImages) > 0 {
			for _, imgInput := range variantReq.ProductImages {
				image := entity.ProductImage{
					ProductID: variantReq.ProductID,
					VariantID: &variant.ID,
					URL:       imgInput.URL,
					IsDefault: imgInput.IsDefault,
				}
				if err := tx.Create(&image).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *VariantRepository) UpdateVariant(ctx context.Context, id string, variantReq *request.VariantPatchRequest) (*entity.Variant, error) {
	// Check if product exists if ProductID is being updated
	var product entity.Product
	if variantReq.ProductID != nil {
		if err := r.db.WithContext(ctx).Preload("Attributes").First(&product, "id = ?", *variantReq.ProductID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg.ProductNotFound
			}
			return nil, err
		}
	} else {
		// Get current variant's product to validate attributes
		var variant entity.Variant
		if err := r.db.WithContext(ctx).First(&variant, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg.VariantNotFound
			}
			return nil, err
		}
		if err := r.db.WithContext(ctx).Preload("Attributes").First(&product, "id = ?", variant.ProductID).Error; err != nil {
			return nil, err
		}
	}

	// Check if attributes exist and belong to product if being updated
	if variantReq.AttributeValues != nil && len(*variantReq.AttributeValues) > 0 {
		attributes := r.CheckAttributeValuesExist(ctx, *variantReq.AttributeValues)
		if attributes == nil {
			return nil, pkg.AttributeValueNotFound
		}
		// Validate that all attribute values belong to product's attributes
		if !r.ValidateAttributeValuesForProduct(ctx, product, attributes) {
			return nil, pkg.InvalidAttributeValueForProduct
		}
	}

	updates := make(map[string]interface{})
	if variantReq.ProductID != nil {
		updates["product_id"] = *variantReq.ProductID
	}
	if variantReq.SKU != nil {
		updates["sku"] = *variantReq.SKU
	}
	if variantReq.BasePrice != nil {
		updates["base_price"] = *variantReq.BasePrice
	}
	if variantReq.ComparePrice != nil {
		updates["compare_price"] = *variantReq.ComparePrice
	}
	if variantReq.Stock != nil {
		updates["stock"] = *variantReq.Stock
	}

	if len(updates) == 0 && variantReq.AttributeValues == nil && variantReq.ProductImages == nil {
		return nil, pkg.NoFieldsToUpdate
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Fetch the variant first
		var variant entity.Variant
		if err := tx.First(&variant, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pkg.VariantNotFound
			}
			return err
		}

		// Update fields if provided
		if len(updates) > 0 {
			result := tx.Model(&variant).Updates(updates)
			if result.Error != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
					return pkg.DuplicateEntry
				}
				return result.Error
			}
		}

		// Update attribute values if provided
		if variantReq.AttributeValues != nil {
			attributes := r.CheckAttributeValuesExist(ctx, *variantReq.AttributeValues)
			variant.AttributeValues = attributes

			// Use FullSaveAssociations to save the association
			if err := tx.Model(&product).Association("AttributeValues").Replace(variant); err != nil {
				return err
			}
		}

		// Update product images if provided
		if variantReq.ProductImages != nil {
			// Delete existing images for this variant
			if err := tx.Where("variant_id = ?", variant.ID).Delete(&entity.ProductImage{}).Error; err != nil {
				return err
			}

			// Add new images
			for _, imgInput := range *variantReq.ProductImages {
				image := entity.ProductImage{
					ProductID: variant.ProductID,
					VariantID: &variant.ID,
					URL:       imgInput.URL,
					IsDefault: imgInput.IsDefault,
				}
				if err := tx.Create(&image).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetVariantByID(ctx, id)
}

func (r *VariantRepository) DeleteVariant(ctx context.Context, id string) error {
	var variant entity.Variant
	if err := r.db.WithContext(ctx).First(&variant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.VariantNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&variant).Association("AttributeValues").Clear(); err != nil {
			return err
		}
		if err := tx.Delete(&variant).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *VariantRepository) CheckAttributeValuesExist(ctx context.Context, attrValueIDs []uint) []entity.AttributeValue {
	var attributeValues []entity.AttributeValue
	if err := r.db.WithContext(ctx).Where("id IN ?", attrValueIDs).Find(&attributeValues).Error; err != nil {
		return nil
	}

	if len(attributeValues) != len(attrValueIDs) {
		return nil
	}

	return attributeValues
}

func (r *VariantRepository) ValidateAttributeValuesForProduct(ctx context.Context, product entity.Product, attrValues []entity.AttributeValue) bool {
	// Get attribute IDs from the product's attributes
	productAttrIDs := make(map[uint]bool)
	for _, attr := range product.Attributes {
		productAttrIDs[attr.ID] = true
	}

	// Check if all attribute values belong to product's attributes
	for _, attrVal := range attrValues {
		if !productAttrIDs[attrVal.AttributeID] {
			return false
		}
	}

	return true
}
