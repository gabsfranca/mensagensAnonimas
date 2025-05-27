package repo

import (
	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"gorm.io/gorm"
)

type GormAdminRepo struct {
	db *gorm.DB
}

func NewGormAdminRepo(db *gorm.DB) *GormAdminRepo {
	return &GormAdminRepo{db}
}

func (r *GormAdminRepo) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *GormAdminRepo) FindByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
