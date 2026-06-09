package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"errors"
	"time"
)

type OrderService interface {
	CreateOrder(userID string, req domain.CreateOrderRequest) (*domain.OrderResponse, error)
	GetByID(id string) (*domain.OrderResponse, error)
	GetByUserID(userID string) ([]domain.OrderResponse, error)
	GetByRestaurantID(restaurantID string) ([]domain.OrderResponse, error)
	GetAvailable() ([]domain.OrderResponse, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	CancelOrder(id string, userID string) error
}

type orderService struct {
	orderRepo      repository.OrderRepository
	restaurantRepo repository.RestaurantRepository
	menuItemRepo   repository.MenuItemRepository
	courierRepo    repository.CourierRepository
	deliveryRepo   repository.DeliveryRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	restaurantRepo repository.RestaurantRepository,
	menuItemRepo repository.MenuItemRepository,
	courierRepo repository.CourierRepository,
	deliveryRepo repository.DeliveryRepository,
) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		restaurantRepo: restaurantRepo,
		menuItemRepo:   menuItemRepo,
		courierRepo:    courierRepo,
		deliveryRepo:   deliveryRepo,
	}
}

func (s *orderService) CreateOrder(userID string, req domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	restaurant, err := s.restaurantRepo.GetByID(req.RestaurantID)
	if err != nil {
		return nil, errors.New("ресторан не найден")
	}
	if restaurant.Status != domain.RestaurantStatusActive {
		return nil, errors.New("ресторан закрыт")
	}

	var totalPrice float64
	var orderItems []domain.OrderItem

	for _, item := range req.Items {
		menuItem, err := s.menuItemRepo.GetByID(item.MenuItemID)
		if err != nil {
			return nil, errors.New("блюдо не найдено")
		}
		if !menuItem.IsAvailable {
			return nil, errors.New("блюдо недоступно: " + menuItem.Name)
		}
		if menuItem.RestaurantID != req.RestaurantID {
			return nil, errors.New("блюдо не принадлежит этому ресторану")
		}

		totalPrice += menuItem.Price * float64(item.Quantity)
		orderItems = append(orderItems, domain.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
			Price:      menuItem.Price,
		})
	}

	order := &domain.Order{
		UserID:        userID,
		RestaurantID:  req.RestaurantID,
		TotalPrice:    totalPrice,
		Status:        domain.OrderStatusPending,
		Address:       req.Address,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: domain.PaymentStatusPending,
		Items:         orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, errors.New("ошибка создания заказа")
	}

	return toOrderResponse(order), nil
}

func (s *orderService) GetByID(id string) (*domain.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("заказ не найден")
	}
	return toOrderResponse(order), nil
}

func (s *orderService) GetByUserID(userID string) ([]domain.OrderResponse, error) {
	orders, err := s.orderRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("ошибка получения заказов")
	}

	var response []domain.OrderResponse
	for _, o := range orders {
		response = append(response, *toOrderResponse(&o))
	}
	return response, nil
}

func (s *orderService) GetByRestaurantID(restaurantID string) ([]domain.OrderResponse, error) {
	orders, err := s.orderRepo.GetByRestaurantID(restaurantID)
	if err != nil {
		return nil, errors.New("ошибка получения заказов")
	}

	var response []domain.OrderResponse
	for _, o := range orders {
		response = append(response, *toOrderResponse(&o))
	}
	return response, nil
}

func (s *orderService) GetAvailable() ([]domain.OrderResponse, error) {
	orders, err := s.orderRepo.GetAvailable()
	if err != nil {
		return nil, errors.New("ошибка получения заказов")
	}

	var response []domain.OrderResponse
	for _, o := range orders {
		response = append(response, *toOrderResponse(&o))
	}
	return response, nil
}

func (s *orderService) UpdateStatus(id string, status domain.OrderStatus) error {
	_, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("заказ не найден")
	}
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *orderService) CancelOrder(id string, userID string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return errors.New("заказ не найден")
	}

	if order.UserID != userID {
		return errors.New("нет прав для отмены заказа")
	}

	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusAccepted {
		return errors.New("заказ нельзя отменить на этом этапе")
	}

	return s.orderRepo.UpdateStatus(id, domain.OrderStatusCancelled)
}

func toOrderResponse(o *domain.Order) *domain.OrderResponse {
	var items []domain.OrderItemResponse
	for _, item := range o.Items {
		items = append(items, domain.OrderItemResponse{
			ID:         item.ID,
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
			Price:      item.Price,
		})
	}

	return &domain.OrderResponse{
		ID:            o.ID,
		UserID:        o.UserID,
		RestaurantID:  o.RestaurantID,
		TotalPrice:    o.TotalPrice,
		Status:        o.Status,
		Address:       o.Address,
		PaymentMethod: o.PaymentMethod,
		PaymentStatus: o.PaymentStatus,
		Items:         items,
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
