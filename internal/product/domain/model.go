package domain

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

type Product struct {
	ID          uuid.UUID       `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:255;not null" json:"name"`
	Description string          `gorm:"size:255" json:"description"`
	Price       float64         `gorm:"type:numeric(10,2);not null" json:"price"`
	Size        json.RawMessage `gorm:"type:jsonb" json:"size"`
	Color       string          `gorm:"size:100" json:"color"`
	Material    *string         `gorm:"size:100" json:"material"`
	ImageURLS   json.RawMessage `gorm:"type:jsonb" json:"imageURLS"`
	IsAvailable bool            `gorm:"not null;default:true" json:"isAvailable"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"deletedAt,omitempty"`
}

type Products []Product

func (p *Product) BeforeCreated() {
	p.ID = uuid.New()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

type FormDataProduct struct {
	Name        string
	Description string
	Price       float64
	Color       string
	IsAvailable bool
	Size        json.RawMessage
	ImageFiles  []*multipart.FileHeader
}
