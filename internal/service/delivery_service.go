package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"errors"
	"time"
)

type DeliveryService interface {
	GetByID(id string) (*domain.DeliveryResponse, error)
	GetByOrderID(orderID string) (*domain.DeliveryResponse, error)
	GetByCourierID(courierID string) ([]domain.DeliveryResponse, error)
	UpdateStatus(id string, status domain.DeliveryStatus) error
}

type deliveryService struct {
	deliveryRepo repository.DeliveryRepository
	orderRepo    repository.OrderRepository
	courierRepo  repository.CourierRepository
}

func NewDeliveryService(
	deliveryRepo repository.DeliveryRepository,
	orderRepo repository.OrderRepository,
	courierRepo repository.CourierRepository,
) DeliveryService {
	return &deliveryService{
		deliveryRepo: deliveryRepo,
		orderRepo:    orderRepo,
		courierRepo:  courierRepo,
	}
}

func (s *deliveryService) GetByID(id string) (*domain.DeliveryResponse, error) {
	delivery, err := s.deliveryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("доставка не найдена")
	}
	return toDeliveryResponse(delivery), nil
}

func (s *deliveryService) GetByOrderID(orderID string) (*domain.DeliveryResponse, error) {
	delivery, err := s.deliveryRepo.GetByOrderID(orderID)
	if err != nil {
		return nil, errors.New("доставка не найдена")
	}
	return toDeliveryResponse(delivery), nil
}

func (s *deliveryService) GetByCourierID(courierID string) ([]domain.DeliveryResponse, error) {
	deliveries, err := s.deliveryRepo.GetByCourierID(courierID)
	if err != nil {
		return nil, errors.New("ошибка получения доставок")
	}

	var response []domain.DeliveryResponse
	for _, d := range deliveries {
		response = append(response, *toDeliveryResponse(&d))
	}
	return response, nil
}

func (s *deliveryService) UpdateStatus(id string, status domain.DeliveryStatus) error {
	delivery, err := s.deliveryRepo.GetByID(id)
	if err != nil {
		return errors.New("доставка не найдена")
	}

	now := time.Now()

	switch status {
	case domain.DeliveryStatusPickedUp:
		s.deliveryRepo.UpdatePickedUpAt(id, &now)
		s.orderRepo.UpdateStatus(delivery.OrderID, domain.OrderStatusPickedUp)
	case domain.DeliveryStatusOnTheWay:
		s.orderRepo.UpdateStatus(delivery.OrderID, domain.OrderStatusOnTheWay)
	case domain.DeliveryStatusDelivered:
		s.deliveryRepo.UpdateDeliveredAt(id, &now)
		s.orderRepo.UpdateStatus(delivery.OrderID, domain.OrderStatusDelivered)
		s.courierRepo.UpdateStatus(delivery.CourierID, domain.CourierStatusOnline)
	case domain.DeliveryStatusFailed:
		s.orderRepo.UpdateStatus(delivery.OrderID, domain.OrderStatusCancelled)
		s.courierRepo.UpdateStatus(delivery.CourierID, domain.CourierStatusOnline)
	}

	return s.deliveryRepo.UpdateStatus(id, status)
}

func toDeliveryResponse(d *domain.Delivery) *domain.DeliveryResponse {
	return &domain.DeliveryResponse{
		ID:          d.ID,
		OrderID:     d.OrderID,
		CourierID:   d.CourierID,
		Status:      d.Status,
		PickedUpAt:  d.PickedUpAt,
		DeliveredAt: d.DeliveredAt,
	}
}
