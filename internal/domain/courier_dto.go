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

type CourierLocationRequest struct {
	Lat float64 `json:"lat" binding:"required"`
	Lng float64 `json:"lng" binding:"required"`
}
