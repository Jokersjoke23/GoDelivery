package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type CourierRepository interface {
	Create(courier *domain.Courier) error
	GetByID(id string) (*domain.Courier, error)
	GetByUserID(userID string) (*domain.Courier, error)
	GetAllOnline() ([]domain.Courier, error)
	GetAll() ([]domain.Courier, error)
	UpdateStatus(id string, status domain.CourierStatus) error
	UpdateLocation(id string, lat float64, lng float64) error
	Delete(id string) error
}

type courierRepository struct {
	db *gorm.DB
}

func NewCourierRepository(db *gorm.DB) CourierRepository {
	return &courierRepository{db: db}
}

func (r *courierRepository) Create(courier *domain.Courier) error {
	return r.db.Create(courier).Error
}

func (r *courierRepository) GetByID(id string) (*domain.Courier, error) {
	var courier domain.Courier
	if err := r.db.Preload("User").First(&courier, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &courier, nil
}

func (r *courierRepository) GetByUserID(userID string) (*domain.Courier, error) {
	var courier domain.Courier
	if err := r.db.Unscoped().First(&courier, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &courier, nil
}

func (r *courierRepository) GetAllOnline() ([]domain.Courier, error) {
	var couriers []domain.Courier
	if err := r.db.Where("status = ?", domain.CourierStatusOnline).Find(&couriers).Error; err != nil {
		return nil, err
	}
	return couriers, nil
}

func (r *courierRepository) GetAll() ([]domain.Courier, error) {
	var couriers []domain.Courier
	if err := r.db.Find(&couriers).Error; err != nil {
		return nil, err
	}
	return couriers, nil
}

func (r *courierRepository) UpdateStatus(id string, status domain.CourierStatus) error {
	return r.db.Model(&domain.Courier{}).Where("id = ?", id).Update("status", status).Error
}

func (r *courierRepository) UpdateLocation(id string, lat float64, lng float64) error {
	return r.db.Model(&domain.Courier{}).Where("id = ?", id).Updates(map[string]interface{}{
		"location_lat": lat,
		"location_lng": lng,
	}).Error
}

func (r *courierRepository) Delete(id string) error {
	return r.db.Unscoped().Delete(&domain.Courier{}, "id = ?", id).Error
}
