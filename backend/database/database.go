package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/models"
	"github.com/go-errors/errors"
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
	var sqlDB *sql.DB // Adicionado para gerenciar a conexão

	maxRetries := 18

	for attempts := 1; attempts <= maxRetries; attempts++ {
		log.Println("tentando conexao com db")
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("tentativa %d falhou: %v", attempts, err)
			time.Sleep(3 * time.Second)
			continue
		}

		// Obter conexão SQL subjacente
		sqlDB, err = db.DB()
		if err != nil {
			log.Printf("falha ao obter DB() %d: %v", attempts, err)
			time.Sleep(3 * time.Second)
			continue
		}

		// Testar conexão com ping
		if err = sqlDB.Ping(); err != nil {
			log.Printf("ping falhou %d: %v", attempts, err)
			sqlDB.Close() // Fechar conexão inválida
			time.Sleep(3 * time.Second)
			continue
		}

		log.Println("banco de dados conectado!")
		return db, nil // Conexão bem-sucedida
	}

	// Tratamento de erro após todas as tentativas
	if err != nil {
		// Usar errors.New se err não for nil
		stackErr := errors.Wrap(err, 1)
		return nil, fmt.Errorf("❌ Nao foi possivel conectar ao db após %d tentativas: %s",
			maxRetries, stackErr.Error())
	}

	return nil, fmt.Errorf("falha desconhecida após %d tentativas", maxRetries)
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("executando migracoes...")
	err := db.AutoMigrate(&models.Report{}, &models.Media{}, &models.Admin{})
	if err != nil {
		// Correção: verificar se o erro não é nil antes de usar Wrap
		stackErr := errors.Wrap(err, 1)
		log.Printf("[ERROR] AutoMigrate falhou: %s\nStacktrace:\n%s",
			err.Error(), stackErr.Error())
	}
	return err
}
