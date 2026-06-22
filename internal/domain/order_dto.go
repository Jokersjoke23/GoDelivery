package domain

type CreateOrderRequest struct {
	RestaurantID  string                   `json:"restaurant_id" binding:"required"`
	Address       string                   `json:"address" binding:"required"`
	PaymentMethod PaymentMethod            `json:"payment_method" binding:"required"`
	Items         []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
}

type CreateOrderItemRequest struct {
	MenuItemID string `json:"menu_item_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required"`
}

type OrderItemResponse struct {
	ID         string  `json:"id"`
	MenuItemID string  `json:"menu_item_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	UserID        string              `json:"user_id"`
	RestaurantID  string              `json:"restaurant_id"`
	TotalPrice    float64             `json:"total_price"`
	DeliveryPrice float64             `json:"delivery_price"`
	GrandTotal    float64             `json:"grand_total"`
	PriceSummary  string              `json:"price_summary"`
	Status        OrderStatus         `json:"status"`
	Address       string              `json:"address"`
	PaymentMethod PaymentMethod       `json:"payment_method"`
	PaymentStatus PaymentStatus       `json:"payment_status"`
	Items         []OrderItemResponse `json:"items"`
}
