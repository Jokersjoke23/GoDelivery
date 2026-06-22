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

// ===== МОКИ =====

type mockPasswordResetRepo struct {
	resets map[string]*domain.PasswordReset
}

func newMockPasswordResetRepo() *mockPasswordResetRepo {
	return &mockPasswordResetRepo{resets: make(map[string]*domain.PasswordReset)}
}

func (m *mockPasswordResetRepo) Create(reset *domain.PasswordReset) error {
	m.resets[reset.Token] = reset
	return nil
}

func (m *mockPasswordResetRepo) GetByToken(token string) (*domain.PasswordReset, error) {
	reset, ok := m.resets[token]
	if !ok {
		return nil, errors.New("не найден")
	}
	return reset, nil
}

func (m *mockPasswordResetRepo) DeleteByUserID(userID string) error {
	for token, r := range m.resets {
		if r.UserID == userID {
			delete(m.resets, token)
		}
	}
	return nil
}

type mockMailer struct{}

func (m *mockMailer) SendResetPassword(to string, token string) error {
	return nil
}

// ===== ХЕЛПЕРЫ =====

func newTestAuthService() AuthService {
	repo := newMockUserRepo()
	resetRepo := newMockPasswordResetRepo()
	h := hasher.NewHasher()
	j := jwt.NewJWT("test-secret", 24*time.Hour)
	return NewAuthService(repo, resetRepo, h, j, &mockMailer{})
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

func TestForgotPassword_Success(t *testing.T) {
	svc := newTestAuthService()

	// сначала регистрируемся
	_, err := svc.Register(domain.CreateUserRequest{
		Name:     "Тест",
		Email:    "test@mail.com",
		Phone:    "+77001234567",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("регистрация не прошла: %v", err)
	}

	// запрашиваем сброс пароля
	err = svc.ForgotPassword("test@mail.com")
	if err != nil {
		t.Fatalf("ожидали успех, получили: %v", err)
	}
}

func TestForgotPassword_UserNotFound(t *testing.T) {
	svc := newTestAuthService()

	err := svc.ForgotPassword("notexist@mail.com")
	if err == nil {
		t.Error("ожидали ошибку — пользователь не найден")
	}
}

func TestResetPassword_Success(t *testing.T) {
	svc := newTestAuthService()

	// регистрируемся
	_, err := svc.Register(domain.CreateUserRequest{
		Name:     "Тест",
		Email:    "test@mail.com",
		Phone:    "+77001234567",
		Password: "oldpassword",
	})
	if err != nil {
		t.Fatalf("регистрация не прошла: %v", err)
	}

	// запрашиваем сброс
	err = svc.ForgotPassword("test@mail.com")
	if err != nil {
		t.Fatalf("forgot password не прошёл: %v", err)
	}

	// берём токен из мока
	resetRepo := newMockPasswordResetRepo()
	authSvc := svc.(*authService)
	token := ""
	for t := range authSvc.passwordResetRepo.(*mockPasswordResetRepo).resets {
		token = t
	}

	// сбрасываем пароль
	err = svc.ResetPassword(token, "newpassword")
	if err != nil {
		t.Fatalf("ожидали успешный сброс пароля: %v", err)
	}

	// проверяем что можем залогиниться с новым паролем
	_, err = svc.Login(domain.LoginRequest{
		Email:    "test@mail.com",
		Password: "newpassword",
	})
	if err != nil {
		t.Fatalf("ожидали успешный логин с новым паролем: %v", err)
	}

	_ = resetRepo
}
