package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
	Audio MediaType = "audio"
)

type Media struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	ReportId  string    `gorm:"type:uuid;not null" json:"reportId"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	Type      MediaType `gorm:"type:varchar(20);not null" json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Media) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return
}
