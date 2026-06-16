package domain

import "time"

type CreateDeliveryRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	CourierID string `json:"courier_id" binding:"required"`
}

type UpdateDeliveryStatusRequest struct {
	Status DeliveryStatus `json:"status" binding:"required"`
}

type DeliveryResponse struct {
	ID            string         `json:"id"`
	OrderID       string         `json:"order_id"`
	CourierID     string         `json:"courier_id"`
	Status        DeliveryStatus `json:"status"`
	DeliveryPrice float64        `'json:"delivery_price"'`
	PickedUpAt    *time.Time     `json:"picked_up_at"`
	DeliveredAt   *time.Time     `json:"delivered_at"`
}
