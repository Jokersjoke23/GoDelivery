package handler

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/service"
	"delivery-app/pkg/response"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	exportService service.ExportService
}

func NewExportHandler(exportService service.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

// @Summary     Создать экспорт
// @Tags        exports
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       input body domain.CreateExportRequest true "Тип экспорта"
// @Success     201 {object} domain.ExportResponse
// @Failure     400 {object} response.ErrorResponse
// @Router      /admin/exports [post]
func (h *ExportHandler) Create(c *gin.Context) {
	var req domain.CreateExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	export, err := h.exportService.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, export)
}

// @Summary     История экспортов
// @Tags        exports
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} []domain.ExportResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /admin/exports [get]
func (h *ExportHandler) GetAll(c *gin.Context) {
	exports, err := h.exportService.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, exports)
}

// @Summary     Статус экспорта
// @Tags        exports
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "ID экспорта"
// @Success     200 {object} domain.ExportResponse
// @Failure     404 {object} response.ErrorResponse
// @Router      /admin/exports/{id} [get]
func (h *ExportHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	export, err := h.exportService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, export)
}

// @Summary     Скачать файл экспорта
// @Tags        exports
// @Produce     octet-stream
// @Security    BearerAuth
// @Param       id path string true "ID экспорта"
// @Success     200
// @Failure     404 {object} response.ErrorResponse
// @Router      /admin/exports/{id}/download [get]
func (h *ExportHandler) Download(c *gin.Context) {
	id := c.Param("id")

	export, err := h.exportService.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	if export.Status != domain.ExportStatusDone {
		response.Error(c, http.StatusBadRequest, "файл ещё не готов")
		return
	}

	if _, err := os.Stat(export.FilePath); os.IsNotExist(err) {
		response.Error(c, http.StatusNotFound, "файл не найден")
		return
	}

	c.File(export.FilePath)
}
