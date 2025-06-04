package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	Recebido  Status = "recebido"
	EmAnalise Status = "em análise"
	Concluido Status = "concluído"
)

type Report struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Message   string    `gorm:"type:text;not null" json:"content"`
	Status    Status    `gorm:"type:varchar(20);default:'recebido'" json:"status"`
	Obs       string    `gorm:"type:text;" json:"obs"`
	Media     []Media   `gorm:"foreignKey:ReportId" json:"media"`
	CreatedAt time.Time `json:"createdAt"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return
}
