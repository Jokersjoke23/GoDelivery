package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary     Мой профиль
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} domain.UserResponse
// @Failure     404 {object} response.ErrorResponse
// @Router      /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, user)
}

// @Summary     Обновить профиль
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       input body domain.UpdateUserRequest true "Данные для обновления"
// @Success     200 {object} domain.UserResponse
// @Failure     400 {object} response.ErrorResponse
// @Router      /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.GetString("userID")

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Update(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, user)
}

// @Summary     Все пользователи (admin)
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} []domain.UserResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /admin/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, users)
}

// @Summary     Удалить пользователя (admin)
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID пользователя"
// @Success     200 {object} map[string]string
// @Failure     400 {object} response.ErrorResponse
// @Router      /admin/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.userService.Delete(id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "пользователь удалён"})
}
