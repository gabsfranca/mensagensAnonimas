package repo

import (
	"context"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"gorm.io/gorm"
)

type MediaRepo interface {
	Create(ctx context.Context, media *models.Media) error
}

type GormMediaRepo struct {
	db *gorm.DB
}

func NewGormMediaRepo(db *gorm.DB) MediaRepo {
	return &GormMediaRepo{db: db}
}

func (r *GormMediaRepo) Create(ctx context.Context, m *models.Media) error {
	return r.db.WithContext(ctx).Create(m).Error
}
