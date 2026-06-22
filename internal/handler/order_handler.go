package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// @Summary     Создать заказ
// @Tags        orders
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       input body domain.CreateOrderRequest true "Данные заказа"
// @Success     201 {object} domain.Order
// @Failure     400 {object} response.ErrorResponse
// @Router      /orders [post]
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

// @Summary     Заказ по ID
// @Tags        orders
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID заказа"
// @Success     200 {object} domain.Order
// @Failure     404 {object} response.ErrorResponse
// @Router      /orders/{id} [get]
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, order)
}

// @Summary     Мои заказы
// @Tags        orders
// @Produce     json
// @Security    BearerAuth
// @Param       page  query int false "Страница" default(1)
// @Param       limit query int false "Лимит"    default(10)
// @Success     200 {object} response.PaginatedResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /orders/my [get]
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := c.GetString("userID")

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, total, err := h.orderService.GetByUserIDPaginated(userID, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, http.StatusOK, orders, page, limit, total)
}

// @Summary     Заказы ресторана
// @Tags        orders
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID ресторана"
// @Success     200 {object} []domain.Order
// @Failure     500 {object} response.ErrorResponse
// @Router      /restaurants/{id}/orders [get]
func (h *OrderHandler) GetByRestaurant(c *gin.Context) {
	restaurantID := c.Param("id")
	orders, err := h.orderService.GetByRestaurantID(restaurantID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, orders)
}

// @Summary     Обновить статус заказа (admin)
// @Tags        orders
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID заказа"
// @Param       input body domain.UpdateOrderStatusRequest true "Новый статус"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /admin/orders/{id}/status [put]
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

// @Summary     Отменить заказ
// @Tags        orders
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID заказа"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /orders/{id}/cancel [post]
func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if err := h.orderService.CancelOrder(id, userID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "заказ отменён"})
}

// @Summary     Все заказы (admin)
// @Tags        orders
// @Produce     json
// @Security    BearerAuth
// @Param       page  query int false "Страница" default(1)
// @Param       limit query int false "Лимит"    default(10)
// @Success     200 {object} response.PaginatedResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /admin/orders [get]
func (h *OrderHandler) GetAll(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, total, err := h.orderService.GetAllPaginated(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, http.StatusOK, orders, page, limit, total)
}
