package repository

import (
	"context"
	"errors"
	"product-service/internal/entity"
	"product-service/internal/request"

	"myshop-shared/pkg"

	"gorm.io/gorm"
)

type ProductImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) *ProductImageRepository {
	return &ProductImageRepository{db: db}
}

func (r *ProductImageRepository) GetAllProductImages(ctx context.Context) ([]*entity.ProductImage, error) {
	var images []*entity.ProductImage

	if err := r.db.WithContext(ctx).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ProductImageRepository) GetProductImageByID(ctx context.Context, id string) (*entity.ProductImage, error) {
	var image = new(entity.ProductImage)
	if err := r.db.WithContext(ctx).First(&image, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ProductImageNotFound
		}
		return nil, err
	}

	return image, nil
}

func (r *ProductImageRepository) GetImagesByProductID(ctx context.Context, productID uint) ([]*entity.ProductImage, error) {
	var images []*entity.ProductImage

	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ProductImageRepository) GetImagesByVariantID(ctx context.Context, variantID uint) ([]*entity.ProductImage, error) {
	var images []*entity.ProductImage

	if err := r.db.WithContext(ctx).Where("variant_id = ?", variantID).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ProductImageRepository) AddProductImage(ctx context.Context, imageReq *request.ProductImageCreateRequest) (*entity.ProductImage, error) {
	// Verify product exists
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", imageReq.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ProductNotFound
		}
		return nil, err
	}

	// Verify variant exists if provided
	if imageReq.VariantID != nil {
		var variant entity.Variant
		if err := r.db.WithContext(ctx).First(&variant, "id = ?", *imageReq.VariantID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, pkg.VariantNotFound
			}
			return nil, err
		}
	}

	image := &entity.ProductImage{
		ProductID: imageReq.ProductID,
		VariantID: imageReq.VariantID,
		URL:       imageReq.URL,
		IsDefault: imageReq.IsDefault,
	}

	if err := r.db.WithContext(ctx).Create(&image).Error; err != nil {
		return nil, err
	}

	return image, nil
}

func (r *ProductImageRepository) UpdateProductImage(ctx context.Context, id string, imageReq *request.ProductImageUpdateRequest) (*entity.ProductImage, error) {
	updates := make(map[string]interface{})

	updates["url"] = imageReq.URL
	updates["is_default"] = imageReq.IsDefault

	if len(updates) == 0 {
		return nil, pkg.NoFieldsToUpdate
	}

	result := r.db.WithContext(ctx).Model(&entity.ProductImage{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, pkg.ProductImageNotFound
	}

	return r.GetProductImageByID(ctx, id)
}

func (r *ProductImageRepository) DeleteProductImage(ctx context.Context, id string) error {
	var image entity.ProductImage
	if err := r.db.WithContext(ctx).First(&image, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ProductImageNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Delete(&image).Error
}

func (r *ProductImageRepository) AddProductImageBatch(ctx context.Context, images []*request.ProductImageCreateRequest) ([]*entity.ProductImage, error) {
	var createdImages []*entity.ProductImage

	for _, imgReq := range images {
		img, err := r.AddProductImage(ctx, imgReq)
		if err != nil {
			return nil, err
		}
		createdImages = append(createdImages, img)
	}

	return createdImages, nil
}

func (r *ProductImageRepository) DeleteProductImagesByProductID(ctx context.Context, productID uint) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&entity.ProductImage{}).Error
}

func (r *ProductImageRepository) DeleteProductImagesByVariantID(ctx context.Context, variantID uint) error {
	return r.db.WithContext(ctx).Where("variant_id = ?", variantID).Delete(&entity.ProductImage{}).Error
}
