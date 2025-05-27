package repo

import "github.com/gabsfranca/mensagensAnonimasRH/models"

type AdminRepo interface {
	Create(admin *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
}
