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
	TakeOrder(courierID string, orderID string) (*domain.DeliveryResponse, error)
	NextStatus(id string, courierID string) error
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

func (s *deliveryService) TakeOrder(courierID string, orderID string) (*domain.DeliveryResponse, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, errors.New("заказ не найден")
	}

	if order.Status != domain.OrderStatusReady {
		return nil, errors.New("заказ недоступен для взятия")
	}

	existing, _ := s.deliveryRepo.GetByOrderID(orderID)
	if existing != nil {
		return nil, errors.New("заказ уже взят другим курьером")
	}

	delivery := &domain.Delivery{
		OrderID:       orderID,
		CourierID:     courierID,
		Status:        domain.DeliveryStatusAssigned,
		DeliveryPrice: 1000,
	}

	if err := s.deliveryRepo.Create(delivery); err != nil {
		return nil, errors.New("ошибка создания доставки")
	}

	s.orderRepo.UpdateStatus(orderID, domain.OrderStatusAccepted)
	s.courierRepo.UpdateStatus(courierID, domain.CourierStatusBusy)

	return toDeliveryResponse(delivery), nil
}

func (s *deliveryService) NextStatus(id string, courierID string) error {
	delivery, err := s.deliveryRepo.GetByID(id)
	if err != nil {
		return errors.New("доставка не найдена")
	}

	if delivery.CourierID != courierID {
		return errors.New("нет прав на эту доставку")
	}

	var nextStatus domain.DeliveryStatus

	switch delivery.Status {
	case domain.DeliveryStatusAssigned:
		nextStatus = domain.DeliveryStatusPickedUp
	case domain.DeliveryStatusPickedUp:
		nextStatus = domain.DeliveryStatusOnTheWay
	case domain.DeliveryStatusOnTheWay:
		nextStatus = domain.DeliveryStatusDelivered
	default:
		return errors.New("доставка уже завершена")
	}

	return s.UpdateStatus(id, nextStatus)
}
