package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type DeliveryRepository interface {
	Create(delivery *domain.Delivery) error
	GetByID(id string) (*domain.Delivery, error)
	GetByOrderID(orderID string) (*domain.Delivery, error)
	GetByCourierID(courierID string) ([]domain.Delivery, error)
	UpdateStatus(id string, status domain.DeliveryStatus) error
}

type deliveryRepository struct {
	db *gorm.DB
}

func NewDeliveryRepository(db *gorm.DB) DeliveryRepository {
	return &deliveryRepository{db: db}
}

func (r *deliveryRepository) Create(delivery *domain.Delivery) error {
	return r.db.Create(delivery).Error
}

func (r *deliveryRepository) GetByID(id string) (*domain.Delivery, error) {
	var delivery domain.Delivery
	if err := r.db.Preload("Order").Preload("Courier").First(&delivery, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *deliveryRepository) GetByOrderID(orderID string) (*domain.Delivery, error) {
	var delivery domain.Delivery
	if err := r.db.First(&delivery, "order_id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *deliveryRepository) GetByCourierID(courierID string) ([]domain.Delivery, error) {
	var deliveries []domain.Delivery
	if err := r.db.Where("courier_id = ?", courierID).Preload("Order").Find(&deliveries).Error; err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *deliveryRepository) UpdateStatus(id string, status domain.DeliveryStatus) error {
	return r.db.Model(&domain.Delivery{}).Where("id = ?", id).Update("status", status).Error
}
