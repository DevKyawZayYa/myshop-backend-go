package main

import (
	"log"
	"math/rand"
	"product-service/config"
	"product-service/internal/entity"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load("../../.env.local"); err != nil {
		log.Println("Warning: .env.local file not found, using environment variables")
	}
}

type DBSeeder struct {
	db *gorm.DB
}

func NewDBSeeder(db *gorm.DB) *DBSeeder {
	return &DBSeeder{db: db}
}

func (s *DBSeeder) SeedCategory() ([]*entity.Category, error) {
	categories := make([]*entity.Category, 11)

	for i := range categories {
		categories[i] = &entity.Category{
			Name:        "Category " + strconv.Itoa(i+1),
			Description: "Description for Category " + strconv.Itoa(i+1),
		}
	}

	if err := s.db.Create(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *DBSeeder) SeedProduct(categories []*entity.Category, attributes []*entity.Attribute) error {

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 50; i++ {
		p := entity.Product{
			Name:         "Product " + strconv.Itoa(i+1),
			Description:  "Description for Product " + strconv.Itoa(i+1),
			BasePrice:    float64((i + 1) * 10),
			ComparePrice: float64((i + 1) * 12),
		}

		if err := tx.Create(&p).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Associate categories with product
		if len(categories) > 0 {
			first := categories[i%len(categories)]
			second := categories[(i+1)%len(categories)]
			if err := tx.Model(&p).Association("Categories").Append(first, second); err != nil {
				tx.Rollback()
				return err
			}
		}

		// Associate random attributes with the product
		if len(attributes) > 0 {
			attrCount := rand.Intn(3) + 2
			selectedAttrs := make([]*entity.Attribute, 0)
			for j := 0; j < attrCount && j < len(attributes); j++ {
				idx := rand.Intn(len(attributes))
				// Avoid duplicates
				isDuplicate := false
				for _, attr := range selectedAttrs {
					if attr.ID == attributes[idx].ID {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					selectedAttrs = append(selectedAttrs, attributes[idx])
				}
			}
			if len(selectedAttrs) > 0 {
				if err := tx.Model(&p).Association("Attributes").Append(selectedAttrs); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		// Add product images
		imageCount := rand.Intn(3) + 1
		for j := 0; j < imageCount; j++ {
			image := entity.ProductImage{
				ProductID: p.ID,
				URL:       "https://via.placeholder.com/400x400?text=Product+" + strconv.Itoa(i+1) + "+Image+" + strconv.Itoa(j+1),
				IsDefault: j == 0,
			}
			if err := tx.Create(&image).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *DBSeeder) SeedAttributes() ([]*entity.Attribute, error) {
	attributes := make([]*entity.Attribute, 10)

	for i := range attributes {
		attributes[i] = &entity.Attribute{
			Name: "Attribute " + strconv.Itoa(i+1),
		}
	}

	if err := s.db.Create(&attributes).Error; err != nil {
		return nil, err
	}

	return attributes, nil
}

func (s *DBSeeder) SeedAttributeValue(attributes []*entity.Attribute) error {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 30; i++ {
		var attrID uint
		if len(attributes) > 0 {
			idx := rand.Intn(len(attributes))
			attrID = attributes[idx].ID
		} else {
			attrID = uint(rand.Intn(10) + 1)
		}

		av := &entity.AttributeValue{
			AttributeID: attrID,
			Value:       "Value " + strconv.Itoa(i+1),
		}
		if err := s.db.Create(&av).Error; err != nil {
			log.Printf("Warning: Failed to create attribute value: %v", err)
			// Continue with next iteration even if this one fails (to handle duplicates)
			continue
		}
	}
	return nil
}

func (s *DBSeeder) SeedAttributeValuesWithReturn(attributes []*entity.Attribute) ([]*entity.AttributeValue, error) {
	rand.Seed(time.Now().UnixNano())
	attributeValues := make([]*entity.AttributeValue, 0)

	for i := 0; i < 30; i++ {
		var attrID uint
		if len(attributes) > 0 {
			idx := rand.Intn(len(attributes))
			attrID = attributes[idx].ID
		} else {
			attrID = uint(rand.Intn(10) + 1)
		}

		av := &entity.AttributeValue{
			AttributeID: attrID,
			Value:       "Value " + strconv.Itoa(i+1),
		}
		if err := s.db.Create(&av).Error; err != nil {
			log.Printf("Warning: Failed to create attribute value: %v", err)
			// Continue with next iteration even if this one fails (to handle duplicates)
			continue
		}
		attributeValues = append(attributeValues, av)
	}
	return attributeValues, nil
}

func (s *DBSeeder) SeedVariants(attributeValues []*entity.AttributeValue) error {
	var products []*entity.Product
	if err := s.db.Find(&products).Error; err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	for _, product := range products {
		// Create 2-3 variants per product
		variantCount := rand.Intn(2) + 2
		for j := 0; j < variantCount; j++ {
			variant := &entity.Variant{
				ProductID:    product.ID,
				SKU:          "SKU-" + strconv.Itoa(int(product.ID)) + "-" + strconv.Itoa(j+1),
				BasePrice:    product.BasePrice,
				ComparePrice: product.BasePrice * 1.2,
				Stock:        rand.Intn(100) + 10,
			}

			if err := s.db.Create(&variant).Error; err != nil {
				log.Printf("Warning: Failed to create variant: %v", err)
				continue
			}

			// Associate random attribute values with the variant
			if len(attributeValues) > 0 {
				attrCount := rand.Intn(3) + 1
				selectedAttrs := make([]*entity.AttributeValue, 0)
				for k := 0; k < attrCount && k < len(attributeValues); k++ {
					idx := rand.Intn(len(attributeValues))
					// Avoid duplicates
					isDuplicate := false
					for _, av := range selectedAttrs {
						if av.ID == attributeValues[idx].ID {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						selectedAttrs = append(selectedAttrs, attributeValues[idx])
					}
				}

				if len(selectedAttrs) > 0 {
					if err := s.db.Model(variant).Association("AttributeValues").Append(selectedAttrs); err != nil {
						log.Printf("Warning: Failed to associate attribute values with variant: %v", err)
					}
				}
			}

			// Add variant images
			imageCount := rand.Intn(2) + 1
			for k := 0; k < imageCount; k++ {
				image := entity.ProductImage{
					ProductID: product.ID,
					VariantID: &variant.ID,
					URL:       "https://via.placeholder.com/400x400?text=Variant+" + strconv.Itoa(int(variant.ID)) + "+Image+" + strconv.Itoa(k+1),
					IsDefault: k == 0,
				}
				if err := s.db.Create(&image).Error; err != nil {
					log.Printf("Warning: Failed to create variant image: %v", err)
				}
			}
		}
	}
	return nil
}

func (s *DBSeeder) Clear() error {
	if err := s.db.Exec("DELETE FROM product_categories").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM variant_attribute_values").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM product_images").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM variants").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM products").Error; err != nil {
		return err
	}

	if err := s.db.Exec("DELETE FROM categories").Error; err != nil {
		return err
	}

	if err := s.db.Exec("DELETE FROM attributes").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM attribute_values").Error; err != nil {
		return err
	}
	return nil
}

func (s *DBSeeder) Seed() error {
	categories, err := s.SeedCategory()
	if err != nil {
		return err
	}

	attributes, err := s.SeedAttributes()
	if err != nil {
		return err
	}

	if err := s.SeedProduct(categories, attributes); err != nil {
		return err
	}

	attributeValues, err := s.SeedAttributeValuesWithReturn(attributes)
	if err != nil {
		return err
	}

	if err := s.SeedVariants(attributeValues); err != nil {
		return err
	}

	return nil
}

func main() {
	db := config.ConnectDB()

	// Auto-migrate the entities
	if err := db.AutoMigrate(
		&entity.Category{},
		&entity.Attribute{},
		&entity.AttributeValue{},
		&entity.Product{},
		&entity.ProductImage{},
		&entity.Variant{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	s := NewDBSeeder(db)

	if err := s.Clear(); err != nil {
		log.Fatalf("Failed to clear database: %v", err)
	}

	if err := s.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully!")
}
