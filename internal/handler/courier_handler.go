package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CourierHandler struct {
	courierService service.CourierService
}

func NewCourierHandler(courierService service.CourierService) *CourierHandler {
	return &CourierHandler{courierService: courierService}
}

func (h *CourierHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.Create(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, courier)
}

func (h *CourierHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, courier)
}

func (h *CourierHandler) GetAllOnline(c *gin.Context) {
	couriers, err := h.courierService.GetAllOnline()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, couriers)
}

func (h *CourierHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	var req struct {
		Status domain.CourierStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.courierService.UpdateStatus(courier.ID, req.Status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "статус обновлён"})
}

func (h *CourierHandler) UpdateLocation(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	var req struct {
		Lat float64 `json:"lat" binding:"required"`
		Lng float64 `json:"lng" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.courierService.UpdateLocation(courier.ID, req.Lat, req.Lng); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "геолокация обновлена"})
}
