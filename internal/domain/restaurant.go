package domain

import (
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
	ID        string           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OwnerID   string           `gorm:"type:uuid;not null" json:"owner_id"`
	Owner     *User            `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Name      string           `gorm:"type:varchar(100);not null" json:"name"`
	Address   string           `gorm:"type:varchar(255);not null" json:"address"`
	Phone     string           `gorm:"type:varchar(20);not null" json:"phone"`
	Status    RestaurantStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
	MenuItems []MenuItem       `gorm:"foreignKey:RestaurantID" json:"menu_items,omitempty"`
	NameRu    string           `gorm:"type:varchar(100)" json:"name_ru"`
	NameEn    string           `gorm:"type:varchar(100)" json:"name_en"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type MenuItem struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	RestaurantID string         `gorm:"type:uuid;not null" json:"restaurant_id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	Price        float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	IsAvailable  bool           `gorm:"default:true" json:"is_available"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	NameRu       string         `gorm:"type:varchar(100)" json:"name_ru"`
	NameEn       string         `gorm:"type:varchar(100)" json:"name_en"`
}

func (MenuItem) TableName() string {
	return "menu_items"
}
