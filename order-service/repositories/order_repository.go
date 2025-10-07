package repositories

import (
	"gorm.io/gorm"
	"order-service/domain/models"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindByProductID(productID int) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	return &order, err
}

func (r *orderRepository) FindByProductID(productID int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("product_id = ?", productID).Find(&orders).Error
	return orders, err
}
