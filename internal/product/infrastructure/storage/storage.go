package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lsendoya/Warewise/internal/product/domain"
	"github.com/lsendoya/Warewise/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type Storage struct {
	db *gorm.DB
}

func New(db *gorm.DB) Storage {
	return Storage{db: db}
}

func (s *Storage) Add(product *domain.Product) (*domain.Product, error) {

	result := s.db.Create(product)

	if result.Error != nil {
		logger.Errorf("%T was not created, s.db.Create() : %v", product, result.Error)
		return &domain.Product{}, result.Error
	}

	logger.Info("Product was created successfully", product)

	return product, nil
}

func (s *Storage) Get(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	result := s.db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (s *Storage) Update(id uuid.UUID, product *domain.Product) (*domain.Product, error) {
	var mdl domain.Product

	updateFields := map[string]interface{}{
		"updatedAt": time.Now(),
	}

	if product.Name != "" {
		updateFields["name"] = product.Name
	}
	if product.Description != "" {
		updateFields["description"] = product.Description
	}

	if product.Size != nil {
		updateFields["size"] = product.Size
	}

	if product.Color != "" {
		updateFields["color"] = product.Color
	}

	if product.Price != 0 {
		updateFields["price"] = product.Price
	}

	if !product.Available {
		updateFields["available"] = false
	} else {
		updateFields["available"] = product.Available
	}

	result := s.db.Model(mdl).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		msgErr := fmt.Sprintf("error updating product with id %s: %v", id, result.Error)
		logger.Errorf(msgErr)
		return nil, result.Error
	}

	logger.Info("product was updated successfully")
	return s.Get(id)

}

func (s *Storage) Delete(id uuid.UUID) error {
	result := s.db.Delete(&domain.Product{}, "id = ?", id)
	if result.Error != nil {
		logger.Errorf("error deleting product with id %s: %v", id, result.Error)
		return result.Error
	}

	logger.Info("product was deleting successfully")

	return nil
}

func (s *Storage) List() (domain.Products, error) {
	var products domain.Products
	result := s.db.Find(&products)

	if result.Error != nil {
		logger.Errorf("error retrieving all products: %v", result.Error)
		return nil, result.Error
	}

	return products, nil
}
