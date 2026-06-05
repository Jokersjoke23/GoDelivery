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

// @Summary     Создать профиль курьера
// @Tags        couriers
// @Produce     json
// @Security    BearerAuth
// @Success     201 {object} domain.Courier
// @Failure     400 {object} response.ErrorResponse
// @Router      /couriers [post]
func (h *CourierHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.Create(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, courier)
}

// @Summary     Мой профиль курьера
// @Tags        couriers
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} domain.Courier
// @Failure     404 {object} response.ErrorResponse
// @Router      /couriers/me [get]
func (h *CourierHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.Success(c, http.StatusOK, courier)
}

// @Summary     Онлайн курьеры (admin)
// @Tags        couriers
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} []domain.Courier
// @Failure     500 {object} response.ErrorResponse
// @Router      /admin/couriers/online [get]
func (h *CourierHandler) GetAllOnline(c *gin.Context) {
	couriers, err := h.courierService.GetAllOnline()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, couriers)
}

// @Summary     Обновить статус курьера
// @Tags        couriers
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       input body domain.CourierStatus true "Новый статус"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /couriers/me/status [put]
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

// @Summary     Обновить геолокацию курьера
// @Tags        couriers
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       input body domain.CourierLocationRequest true "Координаты"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /couriers/me/location [put]
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
