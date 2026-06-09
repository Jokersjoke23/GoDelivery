package domain

import (
	"time"

	"gorm.io/gorm"
)

type ExportType string
type ExportStatus string

const (
	ExportTypeOrders      ExportType = "orders"
	ExportTypeCouriers    ExportType = "couriers"
	ExportTypeRestaurants ExportType = "restaurants"
)

const (
	ExportStatusPending ExportStatus = "pending"
	ExportStatusDone    ExportStatus = "done"
	ExportStatusFailed  ExportStatus = "failed"
)

type Export struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Type      ExportType     `gorm:"type:varchar(20);not null" json:"type"`
	Status    ExportStatus   `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	FilePath  string         `gorm:"type:varchar(255)" json:"file_path"`
	Filters   string         `gorm:"type:text" json:"filters"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Export) TableName() string {
	return "exports"
}

type ExportResponse struct {
	ID        string       `json:"id"`
	Type      ExportType   `json:"type"`
	Status    ExportStatus `json:"status"`
	FilePath  string       `json:"file_path"`
	CreatedAt time.Time    `json:"created_at"`
}

type CreateExportRequest struct {
	Type    ExportType `json:"type" binding:"required"`
	Filters string     `json:"filters,omitempty"`
}
