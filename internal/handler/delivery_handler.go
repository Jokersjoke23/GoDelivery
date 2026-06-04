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

func (h *DeliveryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	delivery, err := h.deliveryService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, delivery)
}

func (h *DeliveryHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("id")
	delivery, err := h.deliveryService.GetByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, delivery)
}

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
