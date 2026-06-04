package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type MenuItemRepository interface {
	Create(item *domain.MenuItem) error
	GetByID(id string) (*domain.MenuItem, error)
	Update(item *domain.MenuItem) error
	Delete(id string) error
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

func (r *menuItemRepository) Update(item *domain.MenuItem) error {
	return r.db.Save(item).Error
}

func (r *menuItemRepository) Delete(id string) error {
	return r.db.Delete(&domain.MenuItem{}, "id = ?", id).Error
}
