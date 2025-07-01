package repo

import (
	"context"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"gorm.io/gorm"
)

type TagRepo interface {
	FindAllTags(ctx context.Context) ([]models.Tag, error)
	RemoveTagFromMessage(Ctx context.Context, messageId string, tagId string) error
	CountReportsByTag(ctx context.Context) ([]models.TagCount, error)
}

type GormTagRepo struct {
	db *gorm.DB
}

func NewGormTagRepo(db *gorm.DB) TagRepo {
	return &GormTagRepo{db: db}
}

func (r *GormTagRepo) FindAllTags(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

func (r *GormTagRepo) RemoveTagFromMessage(ctx context.Context, messageId string, tagId string) error {
	result := r.db.WithContext(ctx).Exec(
		"DELETE FROM report_tags WHERE report_id = ? AND tag_id = ?",
		messageId,
		tagId,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *GormTagRepo) CountReportsByTag(ctx context.Context) ([]models.TagCount, error) {
	var results []models.TagCount

	err := r.db.WithContext(ctx).Raw(`
		SELECT tags.name AS name, COUNT(report_tags.report_id) AS count
		FROM tags
		LEFT JOIN report_tags ON tags.id = report_tags.tag_id
		GROUP BY tags.name
		ORDER BY count DESC
	`).Scan(&results).Error

	return results, err
}
