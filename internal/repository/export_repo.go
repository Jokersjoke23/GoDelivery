package repository

import (
	"delivery-app/internal/domain"

	"gorm.io/gorm"
)

type ExportRepository interface {
	Create(export *domain.Export) error
	GetByID(id string) (*domain.Export, error)
	GetAll() ([]domain.Export, error)
	UpdateStatus(id string, status domain.ExportStatus, filePath string) error
}

type exportRepository struct {
	db *gorm.DB
}

func NewExportRepository(db *gorm.DB) ExportRepository {
	return &exportRepository{db: db}
}

func (r *exportRepository) Create(export *domain.Export) error {
	return r.db.Create(export).Error
}

func (r *exportRepository) GetByID(id string) (*domain.Export, error) {
	var export domain.Export
	if err := r.db.First(&export, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &export, nil
}

func (r *exportRepository) GetAll() ([]domain.Export, error) {
	var exports []domain.Export
	if err := r.db.Order("created_at desc").Find(&exports).Error; err != nil {
		return nil, err
	}
	return exports, nil
}

func (r *exportRepository) UpdateStatus(id string, status domain.ExportStatus, filePath string) error {
	return r.db.Model(&domain.Export{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"file_path": filePath,
	}).Error
}
