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

// GetEnvVar obtém variável de ambiente, com fallback para valor padrão
func GetEnvVar(key string) string {
	value := os.Getenv(key)
	return value
}
