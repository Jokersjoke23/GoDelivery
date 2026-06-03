package domain

import (
	"time"

	"gorm.io/gorm"
)

type Courier struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      string         `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Status      CourierStatus  `gorm:"type:varchar(20);not null;default:'offline'" json:"status"`
	LocationLat float64        `gorm:"type:decimal(10,8)" json:"location_lat"`
	LocationLng float64        `gorm:"type:decimal(11,8)" json:"location_lng"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Courier) TableName() string {
	return "couriers"
}
