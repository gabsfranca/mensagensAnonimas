package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var db *gorm.DB
	var err error

	//vou add retries para que, caso nao consiga conectar o db ele tentar de novo(por conta do docker)

	maxRetries := 18

	for attempts := 1; attempts <= maxRetries; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		fmt.Println("tentando conexao com db")
		if err == nil {
			sqlDB, errPing := db.DB()
			if errPing == nil {
				pingErr := sqlDB.Ping()
				if pingErr != nil {
					log.Println("banco de dados conectado!")
					break
				} else {
					err = pingErr
				}
			}
		}

		log.Printf("tentativa %d de conexão falhou: %s\n", attempts, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("❌ Nao foi possivel conectar ao db após %d tentativas:  %w", maxRetries, err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Report{}, &models.Media{}, &models.Admin{})
}
