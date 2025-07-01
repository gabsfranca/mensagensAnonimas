package models

import (
	"math/rand"
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
	ShortID   string    `gorm:"type:varchar(10);uniqueIndex" json:"shortId"`
	Message   string    `gorm:"type:text;not null" json:"content"`
	Status    Status    `gorm:"type:varchar(20);default:'recebido'" json:"status"`
	Obs       string    `gorm:"type:text;" json:"obs"`
	Media     []Media   `gorm:"foreignKey:ReportId" json:"media"`
	Tags      []Tag     `gorm:"many2many:report_tags;" json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

func generateShortId(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}

	for {
		short := generateShortId(5)
		var count int64
		err := tx.Model(&Report{}).Where("short_id = ?", short).Count(&count).Error
		if err != nil {
			return err
		}

		if count == 0 {
			r.ShortID = short
			break
		}
	}
	return
}
