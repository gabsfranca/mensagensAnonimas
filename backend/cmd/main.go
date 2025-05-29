package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/config"
	"github.com/gabsfranca/mensagensAnonimasRH/router"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)

	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[ERROR] Erro ao criar arquivo de log: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
	} else {
		log.SetOutput(f)
		log.Println("[INFO] Logs sendo salvos em debug.log")
	}

	log.Println("[INFO] Iniciando servidor...")

	err = config.LoadEnvVars()
	if err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Printf("[ERROR] Erro ao carregar .env: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
	}

	dbHost := config.GetEnvVar("DB_HOST")
	dbUser := config.GetEnvVar("DB_USER")
	dbPassword := config.GetEnvVar("DB_PASSWORD")
	dbDatabase := config.GetEnvVar("DB_DATABASE")
	dbPort := config.GetEnvVar("DB_PORT")

	log.Printf(`[INFO] Configurações do banco de dados:
- Host: %s
- Usuário: %s
- Senha: %s
- Banco: %s
- Porta: %s`, dbHost, dbUser, dbPassword, dbDatabase, dbPort)

	log.Println("[INFO] Configurando rotas...")
	r := router.SetupRouter()
	log.Println("[INFO] Rotas configuradas com sucesso!")

	port := config.GetEnvVar("PORT")
	if port == "" {
		port = ":8080"
		log.Println("[WARN] Porta não definida no .env, usando padrão :8080")
	}

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("[INFO] Servidor ouvindo na porta %s", port)
	fmt.Printf("sv ouvindo na porta: ", port)
	if err := server.ListenAndServe(); err != nil {
		stackErr := errors.Wrap(err, 0)
		if stackErr != nil {
			log.Fatalf("[FATAL] Erro ao iniciar servidor: %v\nStacktrace:\n%s", err, stackErr.Stack())
		}
	}
}
