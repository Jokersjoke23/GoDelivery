package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	order, err := h.orderService.CreateOrder(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, order)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetString("userID")
	orders, err := h.orderService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, orders)
}

func (h *OrderHandler) GetByRestaurant(c *gin.Context) {
	restaurantID := c.Param("id")
	orders, err := h.orderService.GetByRestaurantID(restaurantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.orderService.UpdateStatus(id, req.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "статус обновлён"})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if err := h.orderService.CancelOrder(id, userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "заказ отменён"})
}
