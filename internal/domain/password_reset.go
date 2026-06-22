package domain

import "time"

type PasswordReset struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string    `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
