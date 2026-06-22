package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type RestaurantRepository interface {
	Create(restaurant *domain.Restaurant) error
	GetByID(id string) (*domain.Restaurant, error)
	GetAll() ([]domain.Restaurant, error)
	GetByOwnerID(ownerID string) ([]domain.Restaurant, error)
	Update(restaurant *domain.Restaurant) error
	Delete(id string) error
	Search(query string) ([]domain.Restaurant, error)
}

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) Create(restaurant *domain.Restaurant) error {
	return r.db.Create(restaurant).Error
}

func (r *restaurantRepository) GetByID(id string) (*domain.Restaurant, error) {
	var restaurant domain.Restaurant
	if err := r.db.Preload("MenuItems").First(&restaurant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *restaurantRepository) GetAll() ([]domain.Restaurant, error) {
	var restaurants []domain.Restaurant
	if err := r.db.Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *restaurantRepository) GetByOwnerID(ownerID string) ([]domain.Restaurant, error) {
	var restaurants []domain.Restaurant
	if err := r.db.Where("owner_id = ?", ownerID).Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *restaurantRepository) Update(restaurant *domain.Restaurant) error {
	return r.db.Save(restaurant).Error
}

func (r *restaurantRepository) Delete(id string) error {
	return r.db.Delete(&domain.Restaurant{}, "id = ?", id).Error
}

func (r *restaurantRepository) Search(query string) ([]domain.Restaurant, error) {
	var restaurants []domain.Restaurant
	search := "%" + query + "%"
	if err := r.db.Where("name_ru ILIKE ? OR name_en ILIKE ?", search, search).
		Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}
