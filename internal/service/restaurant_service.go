package service

import (
	"delivery-app/internal/domain"
	"delivery-app/internal/repository"
	"errors"

	"github.com/xuri/excelize/v2"
)

type RestaurantService interface {
	Create(ownerID string, req domain.CreateRestaurantRequest) (*domain.RestaurantResponse, error)
	GetByID(id string) (*domain.RestaurantResponse, error)
	GetAll() ([]domain.RestaurantResponse, error)
	GetByOwnerID(ownerID string) ([]domain.RestaurantResponse, error)
	Update(id string, req domain.UpdateRestaurantRequest) (*domain.RestaurantResponse, error)
	Delete(id string) error
	AddMenuItem(restaurantID string, req domain.CreateMenuItemRequest) (*domain.MenuItemResponse, error)
	UpdateMenuItem(id string, req domain.UpdateMenuItemRequest) (*domain.MenuItemResponse, error)
	DeleteMenuItem(id string) error
	GetMenu(restaurantID string) ([]domain.MenuItemResponse, error)
	GetMenuItem(id string) (*domain.MenuItemResponse, error)
	ImportFromExcel(ownerID string, filePath string) (*domain.ImportResult, error)
}

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
	menuItemRepo   repository.MenuItemRepository
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepository, menuItemRepo repository.MenuItemRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo: restaurantRepo,
		menuItemRepo:   menuItemRepo,
	}
}

func (s *restaurantService) Create(ownerID string, req domain.CreateRestaurantRequest) (*domain.RestaurantResponse, error) {
	restaurant := &domain.Restaurant{
		OwnerID: ownerID,
		NameRu:  req.NameRu,
		NameEn:  req.NameEn,
		Address: req.Address,
		Phone:   req.Phone,
		Status:  domain.RestaurantStatusActive,
	}

	if err := s.restaurantRepo.Create(restaurant); err != nil {
		return nil, errors.New("ошибка создания ресторана")
	}

	return toRestaurantResponse(restaurant), nil
}

func (s *restaurantService) GetByID(id string) (*domain.RestaurantResponse, error) {
	restaurant, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("ресторан не найден")
	}
	return toRestaurantResponse(restaurant), nil
}

func (s *restaurantService) GetAll() ([]domain.RestaurantResponse, error) {
	restaurants, err := s.restaurantRepo.GetAll()
	if err != nil {
		return nil, errors.New("ошибка получения ресторанов")
	}

	var response []domain.RestaurantResponse
	for _, r := range restaurants {
		response = append(response, *toRestaurantResponse(&r))
	}
	return response, nil
}

func (s *restaurantService) GetByOwnerID(ownerID string) ([]domain.RestaurantResponse, error) {
	restaurants, err := s.restaurantRepo.GetByOwnerID(ownerID)
	if err != nil {
		return nil, errors.New("ошибка получения ресторанов")
	}

	var response []domain.RestaurantResponse
	for _, r := range restaurants {
		response = append(response, *toRestaurantResponse(&r))
	}
	return response, nil
}

func (s *restaurantService) Update(id string, req domain.UpdateRestaurantRequest) (*domain.RestaurantResponse, error) {
	restaurant, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("ресторан не найден")
	}

	if req.NameRu != nil {
		restaurant.NameRu = *req.NameRu
	}
	if req.NameEn != nil {
		restaurant.NameEn = *req.NameEn
	}
	if req.Address != nil {
		restaurant.Address = *req.Address
	}
	if req.Phone != nil {
		restaurant.Phone = *req.Phone
	}
	if req.Status != nil {
		restaurant.Status = *req.Status
	}

	if err := s.restaurantRepo.Update(restaurant); err != nil {
		return nil, errors.New("ошибка обновления ресторана")
	}

	return toRestaurantResponse(restaurant), nil
}

func (s *restaurantService) Delete(id string) error {
	_, err := s.restaurantRepo.GetByID(id)
	if err != nil {
		return errors.New("ресторан не найден")
	}
	return s.restaurantRepo.Delete(id)
}

func (s *restaurantService) AddMenuItem(restaurantID string, req domain.CreateMenuItemRequest) (*domain.MenuItemResponse, error) {
	_, err := s.restaurantRepo.GetByID(restaurantID)
	if err != nil {
		return nil, errors.New("ресторан не найден")
	}

	item := &domain.MenuItem{
		RestaurantID: restaurantID,
		NameRu:       req.NameRu,
		NameEn:       req.NameEn,
		Price:        req.Price,
		IsAvailable:  req.IsAvailable,
	}

	if err := s.menuItemRepo.Create(item); err != nil {
		return nil, errors.New("ошибка добавления блюда")
	}

	return toMenuItemResponse(item), nil
}

func (s *restaurantService) UpdateMenuItem(id string, req domain.UpdateMenuItemRequest) (*domain.MenuItemResponse, error) {
	item, err := s.menuItemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("блюдо не найдено")
	}

	if req.NameRu != nil {
		item.NameRu = *req.NameRu
	}
	if req.NameEn != nil {
		item.NameEn = *req.NameEn
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}

	if err := s.menuItemRepo.Update(item); err != nil {
		return nil, errors.New("ошибка обновления блюда")
	}

	return toMenuItemResponse(item), nil
}

func (s *restaurantService) DeleteMenuItem(id string) error {
	_, err := s.menuItemRepo.GetByID(id)
	if err != nil {
		return errors.New("блюдо не найдено")
	}
	return s.menuItemRepo.Delete(id)
}

func toRestaurantResponse(r *domain.Restaurant) *domain.RestaurantResponse {
	return &domain.RestaurantResponse{
		ID:      r.ID,
		OwnerID: r.OwnerID,
		NameRu:  r.NameRu,
		NameEn:  r.NameEn,
		Address: r.Address,
		Phone:   r.Phone,
		Status:  r.Status,
	}
}

func toMenuItemResponse(m *domain.MenuItem) *domain.MenuItemResponse {
	return &domain.MenuItemResponse{
		ID:           m.ID,
		RestaurantID: m.RestaurantID,
		NameRu:       m.NameRu,
		NameEn:       m.NameEn,
		Price:        m.Price,
		IsAvailable:  m.IsAvailable,
	}
}

func (s *restaurantService) GetMenu(restaurantID string) ([]domain.MenuItemResponse, error) {
	items, err := s.menuItemRepo.GetByRestaurantID(restaurantID)
	if err != nil {
		return nil, errors.New("ошибка получения меню")
	}

	var response []domain.MenuItemResponse
	for _, item := range items {
		response = append(response, *toMenuItemResponse(&item))
	}
	return response, nil
}

func (s *restaurantService) GetMenuItem(id string) (*domain.MenuItemResponse, error) {
	item, err := s.menuItemRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("блюдо не найдено")
	}
	return toMenuItemResponse(item), nil
}

func (s *restaurantService) ImportFromExcel(ownerID string, filePath string) (*domain.ImportResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, errors.New("ошибка открытия файла")
	}
	defer f.Close()

	rows, err := f.GetRows("Restaurants")
	if err != nil {
		return nil, errors.New("лист Restaurants не найден")
	}

	result := &domain.ImportResult{}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			result.Failed++
			continue
		}

		ownerID := row[1]
		nameEn := row[2]
		address := row[3]
		phone := row[4]
		status := domain.RestaurantStatus(row[5])

		restaurant := &domain.Restaurant{
			OwnerID: ownerID,
			NameEn:  nameEn,
			Address: address,
			Phone:   phone,
			Status:  status,
		}

		if err := s.restaurantRepo.Create(restaurant); err != nil {
			result.Failed++
			continue
		}

		result.Created++
	}

	return result, nil
}
