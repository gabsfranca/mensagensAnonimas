package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Arquivo .env nao encontrado, procurando envvars do docker...")
	}
	return err
}

func GetEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("variavel de ambeinte %s nao definida", key)
	}
	return value
}
