package service

import (
	"errors"
	"log"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/gabsfranca/mensagensAnonimasRH/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo repo.AdminRepo
}

func NewAuthService(repo repo.AdminRepo) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) Register(email, password string) error {
	log.Printf("register chamada")
	existingAdmin, err := s.repo.FindByEmail(email)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("email nao encontrado, criando novo admin")

			hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("erro ao gerar hash no registro: %v", err)
				return err
			}

			admin := &models.Admin{
				Email:    email,
				Password: string(hashed),
			}

			log.Printf("registro feito com sucesso iupiii")

			return s.repo.Create(admin)
		}
		log.Printf("erro ao buscar admin: %v", err)
		return err
	}

	if existingAdmin != nil {
		return errors.New("admin já cadastrado")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash da senha: %v", err)
		return err
	}
	admin := &models.Admin{
		Email:    email,
		Password: string(hashed),
	}
	log.Printf("registro feito com sucesso")
	return s.repo.Create(admin)
}

func (s *AuthService) Login(email, password string) (*models.Admin, error) {
	log.Printf("login chamado")
	admin, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email nao cadastrado")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))

	if err != nil {
		return nil, errors.New("senha incorreta")
	}
	log.Printf("login realiado com sucesso")
	return admin, nil
}
