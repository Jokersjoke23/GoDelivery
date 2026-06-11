package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	deliveryService service.DeliveryService
	courierService  service.CourierService
}

func NewDeliveryHandler(deliveryService service.DeliveryService, courierService service.CourierService) *DeliveryHandler {
	return &DeliveryHandler{
		deliveryService: deliveryService,
		courierService:  courierService,
	}
}

// @Summary     Доставка по ID
// @Tags        deliveries
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID доставки"
// @Success     200 {object} domain.Delivery
// @Failure     404 {object} response.ErrorResponse
// @Router      /deliveries/{id} [get]
func (h *DeliveryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	delivery, err := h.deliveryService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, delivery)
}

// @Summary     Доставка по заказу
// @Tags        deliveries
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID заказа"
// @Success     200 {object} domain.Delivery
// @Failure     404 {object} response.ErrorResponse
// @Router      /orders/{id}/delivery [get]
func (h *DeliveryHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("id")
	delivery, err := h.deliveryService.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, delivery)
}

// @Summary     Мои доставки (курьер)
// @Tags        deliveries
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} []domain.Delivery
// @Failure     500 {object} response.ErrorResponse
// @Router      /couriers/me/deliveries [get]
func (h *DeliveryHandler) GetMyCourierDeliveries(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	deliveries, err := h.deliveryService.GetByCourierID(courier.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, deliveries)
}

// @Summary     Обновить статус доставки
// @Tags        deliveries
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID доставки"
// @Param       input body domain.DeliveryStatus true "Новый статус"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /deliveries/{id}/status [put]
func (h *DeliveryHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status domain.DeliveryStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.deliveryService.UpdateStatus(id, req.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "статус доставки обновлён"})
}

// @Summary     Следующий статус доставки
// @Tags        deliveries
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID доставки"
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Router      /deliveries/{id}/next-status [post]
func (h *DeliveryHandler) NextStatus(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "курьер не найден")
		return
	}

	if err := h.deliveryService.NextStatus(id, courier.ID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "статус обновлён"})
}
