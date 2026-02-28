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

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context, params *params.ProductQueryParam) (
	[]*entity.Product,
	error) {
	var products []*entity.Product

	if err := r.BuildQuery(ctx, params).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) BuildQuery(ctx context.Context, params *params.ProductQueryParam) *gorm.DB {

	baseQuery := r.db.WithContext(ctx).Model(&entity.Product{})

	query := pkg.NewQueryBuilder(baseQuery).
		ApplyPagination(params.Limit, params.Offset).
		ApplyPreload([]string{"Categories", "Attributes"}).Build()

	if params.CategoryID != "" {
		query = query.Joins("JOIN product_categories pc ON pc.product_id = products.id").
			Where("pc.category_id = ?", params.CategoryID)
	}

	return query
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*entity.Product, error) {
	var product = new(entity.Product)
	if err := r.db.WithContext(ctx).Preload("Categories").
		Preload("Attributes").
		First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ProductNotFound
		}
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) AddProduct(ctx context.Context, productReq *request.ProductCreateRequest) error {
	categories := r.CheckCategoriesExist(ctx, productReq.Categories)
	if categories == nil {
		return pkg.CategoryNotFound
	}

	var attributes []entity.Attribute
	if len(productReq.Attributes) > 0 {
		attributes = r.CheckAttributesExist(ctx, productReq.Attributes)
		if attributes == nil {
			return pkg.AttributeNotFound
		}
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		product := entity.Product{
			Name:         productReq.Name,
			Description:  productReq.Description,
			BasePrice:    productReq.BasePrice,
			ComparePrice: productReq.ComparePrice,
		}

		if err := tx.Create(&product).Error; err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return pkg.DuplicateEntry
			}
			return err
		}

		if err := tx.Model(&product).Association("Categories").Append(categories); err != nil {
			return err
		}

		if len(attributes) > 0 {
			if err := tx.Model(&product).Association("Attributes").Append(attributes); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id string, productReq *request.ProductUpdateRequest) (*entity.Product, error) {
	if cats := r.CheckCategoriesExist(ctx, productReq.Categories); cats == nil {
		return nil, pkg.CategoryNotFound
	}

	if len(productReq.Attributes) > 0 {
		if attrs := r.CheckAttributesExist(ctx, productReq.Attributes); attrs == nil {
			return nil, pkg.AttributeNotFound
		}
	}

	updates := make(map[string]interface{})
	updates["name"] = productReq.Name
	updates["description"] = productReq.Description
	updates["base_price"] = productReq.BasePrice
	updates["compare_price"] = productReq.ComparePrice

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			result := tx.Model(&entity.Product{}).Where("id = ?", id).Updates(updates)
			if result.Error != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
					return pkg.DuplicateEntry
				}
				return result.Error
			}
			if result.RowsAffected == 0 {
				return pkg.ProductNotFound
			}
		}

		var product entity.Product
		if err := tx.First(&product, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pkg.ProductNotFound
			}
			return err
		}

		categories := r.CheckCategoriesExist(ctx, productReq.Categories)
		if err := tx.Model(&product).Association("Categories").Replace(categories); err != nil {
			return err
		}

		attributes := r.CheckAttributesExist(ctx, productReq.Attributes)
		if err := tx.Model(&product).Association("Attributes").Replace(attributes); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return r.GetProductByID(ctx, id)
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ProductNotFound
		}
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&product).Association("Categories").Clear(); err != nil {
			return err
		}
		if err := tx.Delete(&product).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProductRepository) CheckCategoriesExist(ctx context.Context, catIDs []uint) []entity.Category {
	var categories []entity.Category
	if err := r.db.WithContext(ctx).Where("id IN ?", catIDs).Find(&categories).Error; err != nil {
		return nil
	}

	if len(categories) != len(catIDs) {
		return nil
	}

	return categories
}

func (r *ProductRepository) CheckAttributesExist(ctx context.Context, attrIDs []uint) []entity.Attribute {
	var attributes []entity.Attribute
	if err := r.db.WithContext(ctx).Where("id IN ?", attrIDs).Find(&attributes).Error; err != nil {
		return nil
	}

	if len(attributes) != len(attrIDs) {
		return nil
	}

	return attributes
}
