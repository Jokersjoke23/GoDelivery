package domain

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID        string         `gorm:"type:uuid;not null;index" json:"user_id"`
	User          *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	RestaurantID  string         `gorm:"type:uuid;not null;index" json:"restaurant_id"`
	Restaurant    *Restaurant    `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"restaurant,omitempty"`
	TotalPrice    float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status        OrderStatus    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Address       string         `gorm:"type:varchar(255);not null" json:"address"`
	PaymentMethod PaymentMethod  `gorm:"type:varchar(20);not null" json:"payment_method"`
	PaymentStatus PaymentStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"payment_status"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID         string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID    string         `gorm:"type:uuid;not null;index" json:"order_id"`
	MenuItemID string         `gorm:"type:uuid;not null" json:"menu_item_id"`
	MenuItem   *MenuItem      `gorm:"foreignKey:MenuItemID" json:"menu_item,omitempty"`
	Quantity   int            `gorm:"not null" json:"quantity"`
	Price      float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
