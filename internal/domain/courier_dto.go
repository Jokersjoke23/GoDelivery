package domain

type CreateCourierRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type UpdateCourierRequest struct {
	Status      *CourierStatus `json:"status,omitempty"`
	LocationLat *float64       `json:"location_lat,omitempty"`
	LocationLng *float64       `json:"location_lng,omitempty"`
}

type CourierResponse struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Status      CourierStatus `json:"status"`
	LocationLat float64       `json:"location_lat"`
	LocationLng float64       `json:"location_lng"`
}
