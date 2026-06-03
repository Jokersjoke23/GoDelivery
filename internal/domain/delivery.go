package domain

import (
	"time"

	"gorm.io/gorm"
)

type Delivery struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID     string         `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	Order       *Order         `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order,omitempty"`
	CourierID   string         `gorm:"type:uuid;not null;index" json:"courier_id"`
	Courier     *Courier       `gorm:"foreignKey:CourierID;constraint:OnDelete:CASCADE" json:"courier,omitempty"`
	Status      DeliveryStatus `gorm:"type:varchar(20);not null;default:'waiting'" json:"status"`
	PickedUpAt  *time.Time     `gorm:"default:null" json:"picked_up_at"`
	DeliveredAt *time.Time     `gorm:"default:null" json:"delivered_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Delivery) TableName() string {
	return "deliveries"
}
