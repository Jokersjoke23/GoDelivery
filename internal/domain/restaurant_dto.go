package domain

type CreateRestaurantRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone" binding:"required,min=10"`
}

type UpdateRestaurantRequest struct {
	Name    *string           `json:"name,omitempty"`
	Address *string           `json:"address,omitempty"`
	Phone   *string           `json:"phone,omitempty"`
	Status  *RestaurantStatus `json:"status,omitempty"`
}

type RestaurantResponse struct {
	ID      string           `json:"id"`
	OwnerID string           `json:"owner_id"`
	Name    string           `json:"name"`
	Address string           `json:"address"`
	Phone   string           `json:"phone"`
	Status  RestaurantStatus `json:"status"`
}

type CreateMenuItemRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Price       float64 `json:"price" binding:"required,min=0"`
	IsAvailable bool    `json:"is_available"`
}

type UpdateMenuItemRequest struct {
	Name        *string  `json:"name,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	IsAvailable *bool    `json:"is_available,omitempty"`
}

type MenuItemResponse struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurant_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
}
