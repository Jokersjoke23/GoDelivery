package repository

import (
	"delivery-app/internal/domain"
	"time"

	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	Create(reset *domain.PasswordReset) error
	GetByToken(token string) (*domain.PasswordReset, error)
	DeleteByUserID(userID string) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(reset *domain.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *passwordResetRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	var reset domain.PasswordReset
	if err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&reset).Error; err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *passwordResetRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.PasswordReset{}).Error
}
