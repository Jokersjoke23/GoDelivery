package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestaurantHandler struct {
	restaurantService service.RestaurantService
}

func NewRestaurantHandler(restaurantService service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: restaurantService}
}

func (h *RestaurantHandler) GetAll(c *gin.Context) {
	restaurants, err := h.restaurantService.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, restaurants)
}

func (h *RestaurantHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	restaurant, err := h.restaurantService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, restaurant)
}

func (h *RestaurantHandler) Create(c *gin.Context) {
	ownerID := c.GetString("userID")
	var req domain.CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	restaurant, err := h.restaurantService.Create(ownerID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, restaurant)
}

func (h *RestaurantHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	restaurant, err := h.restaurantService.Update(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, restaurant)
}

func (h *RestaurantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.restaurantService.Delete(id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "ресторан удалён"})
}

func (h *RestaurantHandler) AddMenuItem(c *gin.Context) {
	restaurantID := c.Param("id")
	var req domain.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.restaurantService.AddMenuItem(restaurantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, item)
}

func (h *RestaurantHandler) UpdateMenuItem(c *gin.Context) {
	itemID := c.Param("itemID")
	var req domain.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.restaurantService.UpdateMenuItem(itemID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, item)
}

func (h *RestaurantHandler) DeleteMenuItem(c *gin.Context) {
	itemID := c.Param("itemID")
	if err := h.restaurantService.DeleteMenuItem(itemID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "блюдо удалено"})
}
