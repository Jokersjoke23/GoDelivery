package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportService interface {
	Create(req domain.CreateExportRequest) (*domain.ExportResponse, error)
	GetAll() ([]domain.ExportResponse, error)
	GetByID(id string) (*domain.ExportResponse, error)
}

type exportService struct {
	exportRepo     repository.ExportRepository
	orderRepo      repository.OrderRepository
	courierRepo    repository.CourierRepository
	restaurantRepo repository.RestaurantRepository
}

func NewExportService(
	exportRepo repository.ExportRepository,
	orderRepo repository.OrderRepository,
	courierRepo repository.CourierRepository,
	restaurantRepo repository.RestaurantRepository,
) ExportService {
	return &exportService{
		exportRepo:     exportRepo,
		orderRepo:      orderRepo,
		courierRepo:    courierRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *exportService) Create(req domain.CreateExportRequest) (*domain.ExportResponse, error) {
	export := &domain.Export{
		Type:    req.Type,
		Status:  domain.ExportStatusPending,
		Filters: req.Filters,
	}

	if err := s.exportRepo.Create(export); err != nil {
		return nil, errors.New("ошибка создания экспорта")
	}

	// запускаем генерацию в фоне — не блокирует сервер
	go s.generate(export.ID, req.Type)

	return toExportResponse(export), nil
}

func (s *exportService) generate(exportID string, exportType domain.ExportType) {
	filePath, err := s.generateFile(exportType)
	if err != nil {
		s.exportRepo.UpdateStatus(exportID, domain.ExportStatusFailed, "")
		return
	}
	s.exportRepo.UpdateStatus(exportID, domain.ExportStatusDone, filePath)
}

func (s *exportService) generateFile(exportType domain.ExportType) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	switch exportType {
	case domain.ExportTypeOrders:
		if err := s.fillOrders(f); err != nil {
			return "", err
		}
	case domain.ExportTypeCouriers:
		if err := s.fillCouriers(f); err != nil {
			return "", err
		}
	case domain.ExportTypeRestaurants:
		if err := s.fillRestaurants(f); err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип экспорта")
	}

	fileName := fmt.Sprintf("exports/%s_%s.xlsx", exportType, time.Now().Format("2006-01-02_15-04-05"))
	if err := f.SaveAs(fileName); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *exportService) fillOrders(f *excelize.File) error {
	sheet := "Orders"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Пользователь", "Ресторан", "Сумма", "Статус", "Адрес", "Оплата", "Дата"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for rowIdx, order := range orders {
		row := rowIdx + 2
		values := []interface{}{
			order.ID,
			order.UserID,
			order.RestaurantID,
			order.TotalPrice,
			string(order.Status),
			order.Address,
			string(order.PaymentMethod),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	return nil
}

func (s *exportService) fillCouriers(f *excelize.File) error {
	sheet := "Couriers"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Пользователь", "Статус", "Широта", "Долгота", "Дата"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	couriers, err := s.courierRepo.GetAll()
	if err != nil {
		return err
	}

	for rowIdx, courier := range couriers {
		row := rowIdx + 2
		values := []interface{}{
			courier.ID,
			courier.UserID,
			string(courier.Status),
			courier.LocationLat,
			courier.LocationLng,
			courier.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	return nil
}

func (s *exportService) fillRestaurants(f *excelize.File) error {
	sheet := "Restaurants"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"ID", "Название", "Адрес", "Телефон", "Статус", "Дата"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	restaurants, err := s.restaurantRepo.GetAll()
	if err != nil {
		return err
	}

	for rowIdx, restaurant := range restaurants {
		row := rowIdx + 2
		values := []interface{}{
			restaurant.ID,
			restaurant.Name,
			restaurant.Address,
			restaurant.Phone,
			string(restaurant.Status),
			restaurant.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	return nil
}

func (s *exportService) GetAll() ([]domain.ExportResponse, error) {
	exports, err := s.exportRepo.GetAll()
	if err != nil {
		return nil, errors.New("ошибка получения экспортов")
	}

	var response []domain.ExportResponse
	for _, e := range exports {
		response = append(response, *toExportResponse(&e))
	}
	return response, nil
}

func (s *exportService) GetByID(id string) (*domain.ExportResponse, error) {
	export, err := s.exportRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("экспорт не найден")
	}
	return toExportResponse(export), nil
}

func toExportResponse(e *domain.Export) *domain.ExportResponse {
	return &domain.ExportResponse{
		ID:        e.ID,
		Type:      e.Type,
		Status:    e.Status,
		FilePath:  e.FilePath,
		CreatedAt: e.CreatedAt,
	}
}
