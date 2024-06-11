package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        string         `json:"id" gorm:"column:id;primaryKey;size:32;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
