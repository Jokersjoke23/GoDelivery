package service

import (
	"delivery-app/internal/domain"
	"testing"
	"time"

	"gorm.io/gorm"
)

// ===== МОКИ =====

type mockCourierRepo struct {
	couriers map[string]*domain.Courier
}

func newMockCourierRepo() *mockCourierRepo {
	return &mockCourierRepo{couriers: make(map[string]*domain.Courier)}
}

func (m *mockCourierRepo) Create(c *domain.Courier) error {
	c.ID = "courier-test-uuid"
	m.couriers[c.UserID] = c
	return nil
}

func (m *mockCourierRepo) GetByID(id string) (*domain.Courier, error) {
	for _, c := range m.couriers {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, errNotFound
}

func (m *mockCourierRepo) GetByUserID(userID string) (*domain.Courier, error) {
	c, ok := m.couriers[userID]
	if !ok {
		return nil, errNotFound
	}
	return c, nil
}

func (m *mockCourierRepo) GetAllOnline() ([]domain.Courier, error) {
	var result []domain.Courier
	for _, c := range m.couriers {
		if c.Status == domain.CourierStatusOnline {
			result = append(result, *c)
		}
	}
	return result, nil
}

func (m *mockCourierRepo) GetAll() ([]domain.Courier, error) {
	var result []domain.Courier
	for _, c := range m.couriers {
		result = append(result, *c)
	}
	return result, nil
}

func (m *mockCourierRepo) UpdateStatus(id string, status domain.CourierStatus) error {
	for _, c := range m.couriers {
		if c.ID == id {
			c.Status = status
		}
	}
	return nil
}

func (m *mockCourierRepo) UpdateLocation(id string, lat float64, lng float64) error {
	for _, c := range m.couriers {
		if c.ID == id {
			c.LocationLat = lat
			c.LocationLng = lng
		}
	}
	return nil
}

func (m *mockCourierRepo) Delete(id string) error {
	for userID, c := range m.couriers {
		if c.ID == id {
			delete(m.couriers, userID)
		}
	}
	return nil
}

// ---

type mockUserRepoForCourier struct {
	users map[string]*domain.User
}

func newMockUserRepoForCourier() *mockUserRepoForCourier {
	m := &mockUserRepoForCourier{users: make(map[string]*domain.User)}
	m.users["user-1"] = &domain.User{
		ID:        "user-1",
		Name:      "Тест",
		Email:     "test@mail.com",
		Role:      domain.UserRoleCustomer,
		IsActive:  true,
		DeletedAt: gorm.DeletedAt{},
	}
	m.users["admin-1"] = &domain.User{
		ID:       "admin-1",
		Name:     "Админ",
		Email:    "admin@mail.com",
		Role:     domain.UserRoleAdmin,
		IsActive: true,
	}
	return m
}

func (m *mockUserRepoForCourier) Create(user *domain.User) error { return nil }
func (m *mockUserRepoForCourier) GetByID(id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}
func (m *mockUserRepoForCourier) GetByEmail(email string) (*domain.User, error) { return nil, nil }
func (m *mockUserRepoForCourier) Update(user *domain.User) error                { return nil }
func (m *mockUserRepoForCourier) Delete(id string) error                        { return nil }
func (m *mockUserRepoForCourier) GetAll() ([]domain.User, error)                { return nil, nil }
func (m *mockUserRepoForCourier) UpdateRole(id string, role domain.UserRole) error {
	if u, ok := m.users[id]; ok {
		u.Role = role
	}
	return nil
}

// ===== ХЕЛПЕР =====

func newTestCourierService() CourierService {
	return NewCourierService(
		newMockCourierRepo(),
		newMockUserRepoForCourier(),
	)
}

// ===== ТЕСТЫ =====

func TestCreateCourier_Success(t *testing.T) {
	svc := newTestCourierService()

	result, err := svc.Create("user-1")

	if err != nil {
		t.Fatalf("ожидали успех, получили: %v", err)
	}
	if result.UserID != "user-1" {
		t.Errorf("ожидали userID user-1, получили %s", result.UserID)
	}
	if result.Status != domain.CourierStatusOffline {
		t.Errorf("ожидали статус offline, получили %s", result.Status)
	}
}

func TestCreateCourier_AdminCannotBeCourier(t *testing.T) {
	svc := newTestCourierService()

	_, err := svc.Create("admin-1")

	if err == nil {
		t.Error("ожидали ошибку — админ не может быть курьером")
	}
}

func TestCreateCourier_AlreadyExists(t *testing.T) {
	svc := newTestCourierService()

	// создаём курьера первый раз
	_, err := svc.Create("user-1")
	if err != nil {
		t.Fatalf("первое создание должно пройти: %v", err)
	}

	// пытаемся создать снова
	_, err = svc.Create("user-1")
	if err == nil {
		t.Error("ожидали ошибку — профиль уже существует")
	}
}

func TestCreateCourier_UserNotFound(t *testing.T) {
	svc := newTestCourierService()

	_, err := svc.Create("nonexistent-user")

	if err == nil {
		t.Error("ожидали ошибку — пользователь не найден")
	}
}

func TestDeleteProfile_Success(t *testing.T) {
	svc := newTestCourierService()

	// создаём курьера
	_, err := svc.Create("user-1")
	if err != nil {
		t.Fatalf("создание не прошло: %v", err)
	}

	// удаляем профиль
	err = svc.DeleteProfile("user-1")
	if err != nil {
		t.Fatalf("ожидали успешное удаление, получили: %v", err)
	}
}

func TestUpdateStatus_Success(t *testing.T) {
	svc := newTestCourierService()

	courier, err := svc.Create("user-1")
	if err != nil {
		t.Fatalf("создание не прошло: %v", err)
	}

	err = svc.UpdateStatus(courier.ID, domain.CourierStatusOnline)
	if err != nil {
		t.Fatalf("ожидали успешное обновление статуса: %v", err)
	}
}

func TestUpdateLocation_Success(t *testing.T) {
	svc := newTestCourierService()

	courier, err := svc.Create("user-1")
	if err != nil {
		t.Fatalf("создание не прошло: %v", err)
	}

	err = svc.UpdateLocation(courier.ID, 43.238949, 76.889709)
	if err != nil {
		t.Fatalf("ожидали успешное обновление локации: %v", err)
	}
}

var _ = time.Now
