package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary     Регистрация
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body domain.CreateUserRequest true "Данные пользователя"
// @Success     201 {object} domain.AuthResponse
// @Failure     400 {object} response.ErrorResponse
// @Router      /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Register(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, result)
}

// @Summary     Вход в систему
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input body domain.LoginRequest true "Email и пароль"
// @Success     200 {object} domain.AuthResponse
// @Failure     401 {object} response.ErrorResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}
