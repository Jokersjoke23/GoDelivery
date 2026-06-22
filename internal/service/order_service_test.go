package service

import (
	"delivery-app/internal/domain"
	"testing"
	"time"
)

// ===== МОКИ =====

type mockOrderRepo struct {
	orders map[string]*domain.Order
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{orders: make(map[string]*domain.Order)}
}

func (m *mockOrderRepo) Create(order *domain.Order) error {
	order.ID = "order-test-uuid"
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepo) GetByID(id string) (*domain.Order, error) {
	order, ok := m.orders[id]
	if !ok {
		return nil, errNotFound
	}
	return order, nil
}

func (m *mockOrderRepo) GetByUserID(userID string) ([]domain.Order, error) {
	var result []domain.Order
	for _, o := range m.orders {
		if o.UserID == userID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (m *mockOrderRepo) GetByRestaurantID(restaurantID string) ([]domain.Order, error) {
	var result []domain.Order
	for _, o := range m.orders {
		if o.RestaurantID == restaurantID {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (m *mockOrderRepo) GetAvailable() ([]domain.Order, error) {
	var result []domain.Order
	for _, o := range m.orders {
		if o.Status == domain.OrderStatusReady {
			result = append(result, *o)
		}
	}
	return result, nil
}

func (m *mockOrderRepo) UpdateStatus(id string, status domain.OrderStatus) error {
	if o, ok := m.orders[id]; ok {
		o.Status = status
	}
	return nil
}

func (m *mockOrderRepo) Delete(id string) error {
	delete(m.orders, id)
	return nil
}

func (m *mockOrderRepo) GetAll() ([]domain.Order, error) {
	var result []domain.Order
	for _, o := range m.orders {
		result = append(result, *o)
	}
	return result, nil
}

// ---

type mockRestaurantRepo struct {
	restaurants map[string]*domain.Restaurant
}

func newMockRestaurantRepo() *mockRestaurantRepo {
	r := &mockRestaurantRepo{restaurants: make(map[string]*domain.Restaurant)}
	r.restaurants["rest-1"] = &domain.Restaurant{
		ID:     "rest-1",
		Status: domain.RestaurantStatusActive,
	}
	return r
}

func (m *mockRestaurantRepo) GetByID(id string) (*domain.Restaurant, error) {
	r, ok := m.restaurants[id]
	if !ok {
		return nil, errNotFound
	}
	return r, nil
}

func (m *mockRestaurantRepo) Create(r *domain.Restaurant) error    { return nil }
func (m *mockRestaurantRepo) GetAll() ([]domain.Restaurant, error) { return nil, nil }
func (m *mockRestaurantRepo) GetByOwnerID(id string) ([]domain.Restaurant, error) {
	return nil, nil
}
func (m *mockRestaurantRepo) Update(r *domain.Restaurant) error { return nil }
func (m *mockRestaurantRepo) Delete(id string) error            { return nil }

// ---

type mockMenuItemRepo struct {
	items map[string]*domain.MenuItem
}

func newMockMenuItemRepo() *mockMenuItemRepo {
	m := &mockMenuItemRepo{items: make(map[string]*domain.MenuItem)}
	m.items["item-1"] = &domain.MenuItem{
		ID:           "item-1",
		RestaurantID: "rest-1",
		NameEn:       "Burger",
		Price:        1200,
		IsAvailable:  true,
	}
	m.items["item-2"] = &domain.MenuItem{
		ID:           "item-2",
		RestaurantID: "rest-1",
		NameEn:       "Pizza",
		Price:        1500,
		IsAvailable:  true,
	}
	return m
}

func (m *mockMenuItemRepo) GetByID(id string) (*domain.MenuItem, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, errNotFound
	}
	return item, nil
}

func (m *mockMenuItemRepo) Create(item *domain.MenuItem) error { return nil }
func (m *mockMenuItemRepo) Update(item *domain.MenuItem) error { return nil }
func (m *mockMenuItemRepo) Delete(id string) error             { return nil }
func (m *mockMenuItemRepo) GetByRestaurantID(restaurantID string) ([]domain.MenuItem, error) {
	return nil, nil
}

// ---

type mockCourierRepoOrder struct{}

func (m *mockCourierRepoOrder) Create(c *domain.Courier) error                 { return nil }
func (m *mockCourierRepoOrder) GetByID(id string) (*domain.Courier, error)     { return nil, nil }
func (m *mockCourierRepoOrder) GetByUserID(id string) (*domain.Courier, error) { return nil, nil }
func (m *mockCourierRepoOrder) GetAllOnline() ([]domain.Courier, error)        { return nil, nil }
func (m *mockCourierRepoOrder) GetAll() ([]domain.Courier, error)              { return nil, nil }
func (m *mockCourierRepoOrder) UpdateStatus(id string, status domain.CourierStatus) error {
	return nil
}
func (m *mockCourierRepoOrder) UpdateLocation(id string, lat float64, lng float64) error {
	return nil
}
func (m *mockCourierRepoOrder) Delete(id string) error { return nil }

// ---

type mockDeliveryRepoOrder struct{}

func (m *mockDeliveryRepoOrder) Create(d *domain.Delivery) error             { return nil }
func (m *mockDeliveryRepoOrder) GetByID(id string) (*domain.Delivery, error) { return nil, nil }
func (m *mockDeliveryRepoOrder) GetByOrderID(id string) (*domain.Delivery, error) {
	return nil, nil
}
func (m *mockDeliveryRepoOrder) GetByCourierID(id string) ([]domain.Delivery, error) {
	return nil, nil
}
func (m *mockDeliveryRepoOrder) UpdateStatus(id string, status domain.DeliveryStatus) error {
	return nil
}
func (m *mockDeliveryRepoOrder) UpdatePickedUpAt(id string, t *time.Time) error  { return nil }
func (m *mockDeliveryRepoOrder) UpdateDeliveredAt(id string, t *time.Time) error { return nil }

// ===== ХЕЛПЕР =====

func newTestOrderService() OrderService {
	return NewOrderService(
		newMockOrderRepo(),
		newMockRestaurantRepo(),
		newMockMenuItemRepo(),
		&mockCourierRepoOrder{},
		&mockDeliveryRepoOrder{},
		nil, // hub не нужен в тестах
	)
}

// ===== ТЕСТЫ =====

func TestCreateOrder_Success(t *testing.T) {
	svc := newTestOrderService()

	req := domain.CreateOrderRequest{
		RestaurantID:  "rest-1",
		Address:       "ул. Абая 10",
		PaymentMethod: domain.PaymentMethodCash,
		Items: []domain.CreateOrderItemRequest{
			{MenuItemID: "item-1", Quantity: 2}, // 1200 * 2 = 2400
			{MenuItemID: "item-2", Quantity: 1}, // 1500 * 1 = 1500
		},
	}

	result, err := svc.CreateOrder("user-1", req)

	if err != nil {
		t.Fatalf("ожидали успех, получили: %v", err)
	}
	if result.TotalPrice != 3900 {
		t.Errorf("ожидали сумму 3900, получили %.2f", result.TotalPrice)
	}
	if result.Status != domain.OrderStatusPending {
		t.Errorf("ожидали статус pending, получили %s", result.Status)
	}
	if len(result.Items) != 2 {
		t.Errorf("ожидали 2 позиции, получили %d", len(result.Items))
	}
}

func TestCreateOrder_RestaurantNotFound(t *testing.T) {
	svc := newTestOrderService()

	req := domain.CreateOrderRequest{
		RestaurantID:  "nonexistent-rest",
		Address:       "ул. Абая 10",
		PaymentMethod: domain.PaymentMethodCash,
		Items: []domain.CreateOrderItemRequest{
			{MenuItemID: "item-1", Quantity: 1},
		},
	}

	_, err := svc.CreateOrder("user-1", req)

	if err == nil {
		t.Error("ожидали ошибку — ресторан не найден")
	}
}

func TestCreateOrder_MenuItemNotAvailable(t *testing.T) {
	svc := newTestOrderService()

	// делаем блюдо недоступным
	menuRepo := newMockMenuItemRepo()
	menuRepo.items["item-1"].IsAvailable = false

	svc = NewOrderService(
		newMockOrderRepo(),
		newMockRestaurantRepo(),
		menuRepo,
		&mockCourierRepoOrder{},
		&mockDeliveryRepoOrder{},
		nil,
	)

	req := domain.CreateOrderRequest{
		RestaurantID:  "rest-1",
		Address:       "ул. Абая 10",
		PaymentMethod: domain.PaymentMethodCash,
		Items: []domain.CreateOrderItemRequest{
			{MenuItemID: "item-1", Quantity: 1},
		},
	}

	_, err := svc.CreateOrder("user-1", req)

	if err == nil {
		t.Error("ожидали ошибку — блюдо недоступно")
	}
}

func TestCancelOrder_Success(t *testing.T) {
	svc := newTestOrderService()

	// создаём заказ
	req := domain.CreateOrderRequest{
		RestaurantID:  "rest-1",
		Address:       "ул. Абая 10",
		PaymentMethod: domain.PaymentMethodCash,
		Items: []domain.CreateOrderItemRequest{
			{MenuItemID: "item-1", Quantity: 1},
		},
	}
	order, _ := svc.CreateOrder("user-1", req)

	// отменяем
	err := svc.CancelOrder(order.ID, "user-1")

	if err != nil {
		t.Fatalf("ожидали успешную отмену, получили: %v", err)
	}
}

func TestCancelOrder_WrongUser(t *testing.T) {
	svc := newTestOrderService()

	req := domain.CreateOrderRequest{
		RestaurantID:  "rest-1",
		Address:       "ул. Абая 10",
		PaymentMethod: domain.PaymentMethodCash,
		Items: []domain.CreateOrderItemRequest{
			{MenuItemID: "item-1", Quantity: 1},
		},
	}
	order, _ := svc.CreateOrder("user-1", req)

	// пытаемся отменить чужой заказ
	err := svc.CancelOrder(order.ID, "user-2")

	if err == nil {
		t.Error("ожидали ошибку — нельзя отменить чужой заказ")
	}
}

func (m *mockOrderRepo) GetAllPaginated(page int, limit int) ([]domain.Order, int, error) {
	var result []domain.Order
	for _, o := range m.orders {
		result = append(result, *o)
	}
	start := (page - 1) * limit
	end := start + limit
	if start > len(result) {
		return []domain.Order{}, len(result), nil
	}
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], len(result), nil
}

func (m *mockOrderRepo) GetByUserIDPaginated(userID string, page int, limit int) ([]domain.Order, int, error) {
	var result []domain.Order
	for _, o := range m.orders {
		if o.UserID == userID {
			result = append(result, *o)
		}
	}
	start := (page - 1) * limit
	end := start + limit
	if start > len(result) {
		return []domain.Order{}, len(result), nil
	}
	if end > len(result) {
		end = len(result)
	}
	return result[start:end], len(result), nil
}

func (m *mockRestaurantRepo) Search(query string) ([]domain.Restaurant, error) {
	return nil, nil
}
func (m *mockMenuItemRepo) Search(query string) ([]domain.MenuItem, error) {
	return nil, nil
}
