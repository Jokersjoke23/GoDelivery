package service

import (
	"delivery-app/internal/domain"
	"delivery-app/pkg/hasher"
	"delivery-app/pkg/jwt"
	"errors"
	"testing"
	"time"
)

// ===== МОКИ =====

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(user *domain.User) error {
	user.ID = "test-uuid-123"
	user.IsActive = true
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("не найден")
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(id string) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("не найден")
}

func (m *mockUserRepo) Update(user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) Delete(id string) error {
	return nil
}

func (m *mockUserRepo) GetAll() ([]domain.User, error) {
	var users []domain.User
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, nil
}

func (m *mockUserRepo) UpdateRole(id string, role domain.UserRole) error {
	for _, u := range m.users {
		if u.ID == id {
			u.Role = role
		}
	}
	return nil
}

// ===== ХЕЛПЕРЫ =====

func newTestAuthService() AuthService {
	repo := newMockUserRepo()
	h := hasher.NewHasher()
	j := jwt.NewJWT("test-secret", 24*time.Hour)
	return NewAuthService(repo, h, j)
}

// ===== ТЕСТЫ =====

func TestRegister_Success(t *testing.T) {
	svc := newTestAuthService()

	req := domain.CreateUserRequest{
		Name:     "Тестовый пользователь",
		Email:    "test@mail.com",
		Phone:    "+77001234567",
		Password: "password123",
	}

	result, err := svc.Register(req)

	if err != nil {
		t.Fatalf("ожидали успех, получили ошибку: %v", err)
	}
	if result.Token == "" {
		t.Error("токен не должен быть пустым")
	}
	if result.User.Email != req.Email {
		t.Errorf("ожидали email %s, получили %s", req.Email, result.User.Email)
	}
	if result.User.Role != domain.UserRoleCustomer {
		t.Errorf("ожидали роль customer, получили %s", result.User.Role)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newTestAuthService()

	req := domain.CreateUserRequest{
		Name:     "Тест",
		Email:    "test@mail.com",
		Phone:    "+77001234567",
		Password: "password123",
	}

	// первая регистрация
	_, err := svc.Register(req)
	if err != nil {
		t.Fatalf("первая регистрация должна пройти успешно: %v", err)
	}

	// вторая регистрация с тем же email
	_, err = svc.Register(req)
	if err == nil {
		t.Error("ожидали ошибку при дублировании email")
	}
}

func TestLogin_Success(t *testing.T) {
	svc := newTestAuthService()

	// сначала регистрируемся
	registerReq := domain.CreateUserRequest{
		Name:     "Тест",
		Email:    "test@mail.com",
		Phone:    "+77001234567",
		Password: "password123",
	}
	_, err := svc.Register(registerReq)
	if err != nil {
		t.Fatalf("регистрация не прошла: %v", err)
	}

	// теперь логинимся
	loginReq := domain.LoginRequest{
		Email:    "test@mail.com",
		Password: "password123",
	}

	result, err := svc.Login(loginReq)

	if err != nil {
		t.Fatalf("ожидали успешный логин, получили: %v", err)
	}
	if result.Token == "" {
		t.Error("токен не должен быть пустым")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newTestAuthService()

	// регистрируемся
	_, err := svc.Register(domain.CreateUserRequest{
		Name:     "Тест",
		Email:    "test@mail.com",
		Phone:    "+77001234567",
		Password: "correctpassword",
	})
	if err != nil {
		t.Fatalf("регистрация не прошла: %v", err)
	}

	// логинимся с неверным паролем
	_, err = svc.Login(domain.LoginRequest{
		Email:    "test@mail.com",
		Password: "wrongpassword",
	})

	if err == nil {
		t.Error("ожидали ошибку при неверном пароле")
	}
}

func TestLogin_NotFound(t *testing.T) {
	svc := newTestAuthService()

	_, err := svc.Login(domain.LoginRequest{
		Email:    "notexist@mail.com",
		Password: "password123",
	})

	if err == nil {
		t.Error("ожидали ошибку — пользователь не существует")
	}
}
