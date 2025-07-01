package service

import (
	"context"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
)

type TagService interface {
	GetAvailableTags(ctx context.Context) ([]models.Tag, error)
	RemoveTagFromMessage(ctx context.Context, messageId string, tagId string) error
	CountReportsByTag(ctx context.Context) ([]models.TagCount, error)
}

type tagService struct {
	repo repo.TagRepo
}

func NewTagService(r repo.TagRepo) TagService {
	return &tagService{repo: r}
}

func (s *tagService) GetAvailableTags(ctx context.Context) ([]models.Tag, error) {
	return s.repo.FindAllTags(ctx)
}

func (s *tagService) RemoveTagFromMessage(ctx context.Context, messageId string, tagId string) error {
	return s.repo.RemoveTagFromMessage(ctx, messageId, tagId)
}

func (s *tagService) CountReportsByTag(ctx context.Context) ([]models.TagCount, error) {
	return s.repo.CountReportsByTag(ctx)
}
