package domain

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserEmail string         `gorm:"size:255;not null;primaryKey" json:"userEmail"`
	UserUUID  string         `gorm:"size:100;not null" json:"userUUID"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Birthdate string         `gorm:"size:50;not null" json:"birthdate"`
	IsAdmin   bool           `json:"isAdmin"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
