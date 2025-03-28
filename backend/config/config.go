package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Erro ao carregar variaveis de ambiente")
	}
	return nil
}

func GetEnvVar(key string) string {
	return os.Getenv(key)
}
