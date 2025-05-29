package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-errors/errors"
	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	err := godotenv.Load()
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		log.Printf("[WARN] Arquivo .env não encontrado, tentando usar variáveis de ambiente do sistema.\nStacktrace:\n%s", stackErr.Stack())
		// Continua execução, apenas avisa
		return fmt.Errorf("arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
	log.Println("[INFO] Variáveis de ambiente carregadas com sucesso")
	return nil
}

func GetEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		stackErr := errors.New(fmt.Sprintf("Variável de ambiente %s não definida", key))
		log.Printf("[WARN] Variável de ambiente '%s' não definida.\nStacktrace:\n%s", key, stackErr.Stack())
	}
	return value
}
