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

// @Summary     Все рестораны
// @Tags        restaurants
// @Produce     json
// @Success     200 {object} []domain.Restaurant
// @Failure     500 {object} response.ErrorResponse
// @Router      /restaurants [get]
func (h *RestaurantHandler) GetAll(c *gin.Context) {
	restaurants, err := h.restaurantService.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, restaurants)
}

// @Summary     Ресторан по ID
// @Tags        restaurants
// @Produce     json
// @Param       id path string true "ID ресторана"
// @Success     200 {object} domain.Restaurant
// @Failure     404 {object} response.ErrorResponse
// @Router      /restaurants/{id} [get]
func (h *RestaurantHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	restaurant, err := h.restaurantService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, restaurant)
}

// @Summary     Создать ресторан
// @Tags        restaurants
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       input body domain.CreateRestaurantRequest true "Данные ресторана"
// @Success     201 {object} domain.Restaurant
// @Failure     400 {object} response.ErrorResponse
// @Router      /restaurants [post]
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

// @Summary     Обновить ресторан
// @Tags        restaurants
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID ресторана"
// @Param       input body domain.UpdateRestaurantRequest true "Данные для обновления"
// @Success     200 {object} domain.Restaurant
// @Failure     400 {object} response.ErrorResponse
// @Router      /restaurants/{id} [put]
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

// @Summary     Удалить ресторан
// @Tags        restaurants
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID ресторана"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /restaurants/{id} [delete]
func (h *RestaurantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.restaurantService.Delete(id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "ресторан удалён"})
}

// @Summary     Добавить блюдо в меню
// @Tags        restaurants
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID ресторана"
// @Param       input body domain.CreateMenuItemRequest true "Данные блюда"
// @Success     201 {object} domain.MenuItem
// @Failure     400 {object} response.ErrorResponse
// @Router      /restaurants/{id}/menu [post]
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

// @Summary     Обновить блюдо
// @Tags        restaurants
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       itemID path string true "ID блюда"
// @Param       input body domain.UpdateMenuItemRequest true "Данные для обновления"
// @Success     200 {object} domain.MenuItem
// @Failure     400 {object} response.ErrorResponse
// @Router      /restaurants/menu/{itemID} [put]
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

// @Summary     Удалить блюдо
// @Tags        restaurants
// @Produce     json
// @Security    BearerAuth
// @Param       itemID path string true "ID блюда"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /restaurants/menu/{itemID} [delete]
func (h *RestaurantHandler) DeleteMenuItem(c *gin.Context) {
	itemID := c.Param("itemID")
	if err := h.restaurantService.DeleteMenuItem(itemID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "блюдо удалено"})
}
