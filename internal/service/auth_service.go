package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"delivery-app/pkg/hasher"
	"delivery-app/pkg/jwt"
	"delivery-app/pkg/mailer"
	"errors"
	"fmt"
	"time"
)

type AuthService interface {
	Register(req domain.CreateUserRequest) (*domain.AuthResponse, error)
	Login(req domain.LoginRequest) (*domain.AuthResponse, error)
	ForgotPassword(email string) error
	ResetPassword(token string, newPassword string) error
}

type authService struct {
	userRepo          repository.UserRepository
	passwordResetRepo repository.PasswordResetRepository
	hasher            hasher.Hasher
	jwt               jwt.JWT
	mailer            mailer.Mailer
}

func NewAuthService(
	userRepo repository.UserRepository,
	passwordResetRepo repository.PasswordResetRepository,
	hasher hasher.Hasher,
	jwt jwt.JWT,
	mailer mailer.Mailer,
) AuthService {
	return &authService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		hasher:            hasher,
		jwt:               jwt,
		mailer:            mailer,
	}
}

func (s *authService) Register(req domain.CreateUserRequest) (*domain.AuthResponse, error) {
	existing, _ := s.userRepo.GetByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     domain.UserRoleCustomer,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("ошибка создания пользователя")
	}

	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *authService) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	if !s.hasher.Compare(req.Password, user.Password) {
		return nil, errors.New("неверный email или пароль")
	}

	if !user.IsActive {
		return nil, errors.New("аккаунт заблокирован")
	}

	token, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("ошибка генерации токена")
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	s.passwordResetRepo.DeleteByUserID(user.ID)

	token := fmt.Sprintf("%d", time.Now().UnixNano())

	reset := &domain.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.passwordResetRepo.Create(reset); err != nil {
		return errors.New("ошибка создания токена")
	}

	return s.mailer.SendResetPassword(email, token)
}

func (s *authService) ResetPassword(token string, newPassword string) error {
	reset, err := s.passwordResetRepo.GetByToken(token)
	if err != nil {
		return errors.New("токен недействителен или истёк")
	}

	hashed, err := s.hasher.Hash(newPassword)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	user, err := s.userRepo.GetByID(reset.UserID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	user.Password = hashed
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка обновления пароля")
	}

	return s.passwordResetRepo.DeleteByUserID(reset.UserID)
}
