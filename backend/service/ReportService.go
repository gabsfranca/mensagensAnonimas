package service

import (
	"context"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
)

type ReportService interface {
	GetAll(ctx context.Context) ([]models.Report, error)
	GetByID(ctx context.Context, id string) (*models.Report, error)
	ChangeStatus(ctx context.Context, id string, status models.Status) error
	AddObs(ctx context.Context, id string, obs string) error
	GetObs(ctx context.Context, id string) (string, error)
	AddTags(ctx context.Context, reportID string, tagIDs []string) error
}

type reportService struct {
	repo repo.ReportRepo
}

func NewReportService(r repo.ReportRepo) ReportService {
	return &reportService{repo: r}
}

func (s *reportService) GetAll(ctx context.Context) ([]models.Report, error) {
	return s.repo.FindAll(ctx)
}

func (s *reportService) GetByID(ctx context.Context, id string) (*models.Report, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *reportService) GetObs(ctx context.Context, id string) (string, error) {
	return s.repo.GetObsById(ctx, id)
}

func (s *reportService) ChangeStatus(ctx context.Context, id string, status models.Status) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *reportService) AddObs(ctx context.Context, id string, obs string) error {
	return s.repo.UpdateObs(ctx, id, obs)
}

func (s *reportService) AddTags(ctx context.Context, reportID string, tagIDs []string) error {
	return s.repo.AddTags(ctx, reportID, tagIDs)
}
