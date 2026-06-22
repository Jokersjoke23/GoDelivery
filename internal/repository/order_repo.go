package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	GetByUserID(userID string) ([]domain.Order, error)
	GetByRestaurantID(restaurantID string) ([]domain.Order, error)
	GetAvailable() ([]domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	Delete(id string) error
	GetAll() ([]domain.Order, error)
	GetAllPaginated(page int, limit int) ([]domain.Order, int, error)
	GetByUserIDPaginated(userID string, page int, limit int) ([]domain.Order, int, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.Preload("Items").Preload("Items.MenuItem").Preload("Delivery").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(userID string) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Where("user_id = ?", userID).Preload("Items").Preload("Delivery").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetByRestaurantID(restaurantID string) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Where("restaurant_id = ?", restaurantID).Preload("Items").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Delete(id string) error {
	return r.db.Delete(&domain.Order{}, "id = ?", id).Error
}

func (r *orderRepository) GetAvailable() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Where("status = ?", domain.OrderStatusReady).
		Preload("Items").
		Preload("Restaurant").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.
		Preload("User").
		Preload("Restaurant").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetAllPaginated(page int, limit int) ([]domain.Order, int, error) {
	var orders []domain.Order
	var total int64

	r.db.Model(&domain.Order{}).Count(&total)

	if err := r.db.
		Preload("User").
		Preload("Restaurant").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (r *orderRepository) GetByUserIDPaginated(userID string, page int, limit int) ([]domain.Order, int, error) {
	var orders []domain.Order
	var total int64

	r.db.Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total)

	if err := r.db.
		Where("user_id = ?", userID).
		Preload("Items").
		Preload("Delivery").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}
