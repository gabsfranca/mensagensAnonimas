package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"gorm.io/gorm"
)

type ReportRepo interface {
	Create(ctx context.Context, report *models.Report) error
	FindAll(ctx context.Context) ([]models.Report, error)
	FindByID(ctx context.Context, id string) (*models.Report, error)
	GetObsById(ctx context.Context, id string) (string, error)
	UpdateStatus(ctx context.Context, id string, status models.Status) error
	UpdateObs(ctx context.Context, id string, obs string) error
	FindByIdWithMedia(ctx context.Context, id string) (*models.Report, error)
	AddTags(ctx context.Context, reportID string, tagIDs []string) error
}

type GormReportRepo struct {
	db *gorm.DB
}

func NewGormReportRepo(db *gorm.DB) ReportRepo {
	return &GormReportRepo{db: db}
}

func (r *GormReportRepo) Create(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *GormReportRepo) FindAll(ctx context.Context) ([]models.Report, error) {
	var reports []models.Report
	err := r.db.WithContext(ctx).
		Preload("Media").
		Order("created_at desc").
		Find(&reports).Error
	return reports, err
}

func (r *GormReportRepo) FindByID(ctx context.Context, id string) (*models.Report, error) {
	var report models.Report
	err := r.db.WithContext(ctx).
		Preload("Media").
		Preload("Tags").
		First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *GormReportRepo) GetObsById(ctx context.Context, id string) (string, error) {
	var obs string

	err := r.db.WithContext(ctx).
		Model(&models.Report{}).
		Select("obs").
		Where("id = ?", id).
		First(&obs).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("Sem observacoes")
		}
		return "", fmt.Errorf("erro ao buscar observacao: %v", err)
	}

	return obs, nil
}

func (r *GormReportRepo) UpdateStatus(ctx context.Context, id string, status models.Status) error {
	report := models.Report{ID: id}
	return r.db.WithContext(ctx).
		Model(&report).
		Update("status", status).
		Error
}

func (r *GormReportRepo) UpdateObs(ctx context.Context, id string, obs string) error {
	report := models.Report{ID: id}
	return r.db.WithContext(ctx).
		Model(&report).
		Update("obs", obs).
		Error
}

func (r *GormReportRepo) FindByIdWithMedia(ctx context.Context, id string) (*models.Report, error) {
	var report models.Report

	err := r.db.WithContext(ctx).
		Preload("Media").
		First(&report, "id = ?", id).Error
	return &report, err
}

func (r *GormReportRepo) AddTags(ctx context.Context, reportID string, tagIDs []string) error {
	report := models.Report{ID: reportID}
	tags := []models.Tag{}
	for _, tagID := range tagIDs {
		var tag models.Tag
		if err := r.db.WithContext(ctx).First(&tag, "id = ?", tagID).Error; err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	return r.db.WithContext(ctx).Model(&report).Association("Tags").Append(tags)
}
