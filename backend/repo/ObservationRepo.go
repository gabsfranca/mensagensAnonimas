package repo

import (
	"context"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"gorm.io/gorm"
)

type ObservationRepo interface {
	Create(ctx context.Context, obs *models.Observation) error
	FindByReportId(ctx context.Context, reportId string) ([]models.Observation, error)
	FindByShortId(ctx context.Context, shortId string) (*models.Report, error)
}

type GormObservationRepo struct {
	db *gorm.DB
}

func NewGormObservationRepo(db *gorm.DB) ObservationRepo {
	return &GormObservationRepo{db: db}
}

func (r *GormObservationRepo) Create(ctx context.Context, obs *models.Observation) error {
	return r.db.WithContext(ctx).Create(obs).Error
}

func (r *GormObservationRepo) FindByReportId(ctx context.Context, reportId string) ([]models.Observation, error) {
	var observations []models.Observation

	err := r.db.WithContext(ctx).
		Where("report_id = ?", reportId).
		Order("created_at asc").
		Find(&observations).Error

	return observations, err
}

func (r *GormObservationRepo) FindByShortId(ctx context.Context, shortId string) (*models.Report, error) {
	var report models.Report
	err := r.db.WithContext(ctx).Where("short_id = ?", shortId).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
