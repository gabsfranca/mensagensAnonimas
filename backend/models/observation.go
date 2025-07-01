package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Observation struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	ReportID  string    `gorm:"type:uuid;not null" json:"reportId"`
	Author    string    `gorm:"type:varchar(20);not null" json:"author"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (o *Observation) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.NewString()
	o.CreatedAt = time.Now()
	return
}
