package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"delivery-app/pkg/hasher"
	"delivery-app/pkg/jwt"
	"errors"
)

type AuthService interface {
	Register(req domain.CreateUserRequest) (*domain.AuthResponse, error)
	Login(req domain.LoginRequest) (*domain.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	hasher   hasher.Hasher
	jwt      jwt.JWT
}

func NewAuthService(userRepo repository.UserRepository, hasher hasher.Hasher, jwt jwt.JWT) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
		jwt:      jwt,
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
