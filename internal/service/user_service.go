package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"errors"
)

type UserService interface {
	GetByID(id string) (*domain.UserResponse, error)
	Update(id string, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	Delete(id string) error
	GetAll() ([]domain.UserResponse, error)
	AssignRole(id string, role domain.UserRole) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	return &domain.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}

func (s *userService) Update(id string, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("ошибка обновления пользователя")
	}

	return &domain.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}

func (s *userService) Delete(id string) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("пользователь не найден")
	}
	return s.userRepo.Delete(id)
}

func (s *userService) GetAll() ([]domain.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, errors.New("ошибка получения пользователей")
	}

	var response []domain.UserResponse
	for _, user := range users {
		response = append(response, domain.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			IsActive: user.IsActive,
		})
	}
	return response, nil
}

func (s *userService) AssignRole(id string, role domain.UserRole) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if user.Role == domain.UserRoleAdmin {
		return errors.New("нельзя изменить роль администратора")
	}

	return s.userRepo.UpdateRole(id, role)
}
