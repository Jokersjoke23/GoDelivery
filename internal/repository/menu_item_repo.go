package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type MenuItemRepository interface {
	Create(item *domain.MenuItem) error
	GetByID(id string) (*domain.MenuItem, error)
	GetByRestaurantID(restaurantID string) ([]domain.MenuItem, error)
	Update(item *domain.MenuItem) error
	Delete(id string) error
	Search(query string) ([]domain.MenuItem, error)
}

type menuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) Create(item *domain.MenuItem) error {
	return r.db.Create(item).Error
}

func (r *menuItemRepository) GetByID(id string) (*domain.MenuItem, error) {
	var item domain.MenuItem
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *menuItemRepository) GetByRestaurantID(restaurantID string) ([]domain.MenuItem, error) {
	var items []domain.MenuItem
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *menuItemRepository) Update(item *domain.MenuItem) error {
	return r.db.Save(item).Error
}

func (r *menuItemRepository) Delete(id string) error {
	return r.db.Delete(&domain.MenuItem{}, "id = ?", id).Error
}

func (r *menuItemRepository) Search(query string) ([]domain.MenuItem, error) {
	var items []domain.MenuItem
	search := "%" + query + "%"
	if err := r.db.Where("name_ru ILIKE ? OR name_en ILIKE ?", search, search).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
