package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"errors"
)

type CourierService interface {
	Create(userID string) (*domain.CourierResponse, error)
	GetByID(id string) (*domain.CourierResponse, error)
	GetByUserID(userID string) (*domain.CourierResponse, error)
	GetAllOnline() ([]domain.CourierResponse, error)
	UpdateStatus(id string, status domain.CourierStatus) error
	UpdateLocation(id string, lat float64, lng float64) error
	Delete(id string) error
	DeleteProfile(userID string) error
}

type courierService struct {
	courierRepo repository.CourierRepository
	userRepo    repository.UserRepository
}

func NewCourierService(courierRepo repository.CourierRepository, userRepo repository.UserRepository) CourierService {
	return &courierService{
		courierRepo: courierRepo,
		userRepo:    userRepo,
	}
}

func (s *courierService) Create(userID string) (*domain.CourierResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}

	if user.Role == domain.UserRoleAdmin {
		return nil, errors.New("администратор не может быть курьером")
	}

	existing, _ := s.courierRepo.GetByUserID(userID)
	if existing != nil && !existing.DeletedAt.Valid {
		return nil, errors.New("профиль курьера уже существует")
	}

	if err := s.userRepo.UpdateRole(userID, domain.UserRoleCourier); err != nil {
		return nil, errors.New("ошибка обновления роли")
	}

	courier := &domain.Courier{
		UserID: userID,
		Status: domain.CourierStatusOffline,
	}

	if err := s.courierRepo.Create(courier); err != nil {
		return nil, errors.New("ошибка создания профиля курьера")
	}

	return toCourierResponse(courier), nil
}

func (s *courierService) GetByID(id string) (*domain.CourierResponse, error) {
	courier, err := s.courierRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("курьер не найден")
	}
	return toCourierResponse(courier), nil
}

func (s *courierService) GetByUserID(userID string) (*domain.CourierResponse, error) {
	courier, err := s.courierRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("курьер не найден")
	}
	return toCourierResponse(courier), nil
}

func (s *courierService) GetAllOnline() ([]domain.CourierResponse, error) {
	couriers, err := s.courierRepo.GetAllOnline()
	if err != nil {
		return nil, errors.New("ошибка получения курьеров")
	}

	var response []domain.CourierResponse
	for _, c := range couriers {
		response = append(response, *toCourierResponse(&c))
	}
	return response, nil
}

func (s *courierService) UpdateStatus(id string, status domain.CourierStatus) error {
	_, err := s.courierRepo.GetByID(id)
	if err != nil {
		return errors.New("курьер не найден")
	}
	return s.courierRepo.UpdateStatus(id, status)
}

func (s *courierService) UpdateLocation(id string, lat float64, lng float64) error {
	_, err := s.courierRepo.GetByID(id)
	if err != nil {
		return errors.New("курьер не найден")
	}
	return s.courierRepo.UpdateLocation(id, lat, lng)
}

func (s *courierService) Delete(id string) error {
	_, err := s.courierRepo.GetByID(id)
	if err != nil {
		return errors.New("курьер не найден")
	}
	return s.courierRepo.Delete(id)
}

func (s *courierService) DeleteProfile(userID string) error {
	courier, err := s.courierRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("профиль курьера не найден")
	}

	if err := s.courierRepo.Delete(courier.ID); err != nil {
		return errors.New("ошибка удаления профиля курьера")
	}

	if err := s.userRepo.UpdateRole(userID, domain.UserRoleCustomer); err != nil {
		return errors.New("ошибка обновления роли")
	}

	return nil
}

func toCourierResponse(c *domain.Courier) *domain.CourierResponse {
	return &domain.CourierResponse{
		ID:          c.ID,
		UserID:      c.UserID,
		Status:      c.Status,
		LocationLat: c.LocationLat,
		LocationLng: c.LocationLng,
	}
}
