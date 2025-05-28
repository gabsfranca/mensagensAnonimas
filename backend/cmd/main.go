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
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := config.LoadEnvVars()
	if err != nil {
		log.Println("Erro .env: ", err)
	}
	dbHost := config.GetEnvVar("DB_HOST")
	dbUser := config.GetEnvVar("DB_USER")
	dbPassword := config.GetEnvVar("DB_PASSWORD")
	dbDatabase := config.GetEnvVar("DB_DATABASE")
	dbPort := config.GetEnvVar("DB_PORT")

	log.Printf("Host banco de dados: %s\nusuário banco de dados: %s, \nsenha banco de dados: %s, \nnome banco de dados: %s, \nporta banco de dados: %s", dbHost, dbUser, dbPassword, dbDatabase, dbPort)

	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)

	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Erro ao criar arquivo de log: %v", err)
	} else {
		log.SetOutput(f)
		fmt.Println("logs sendo salvos em debug.log")
	}

	fmt.Println("iniciandop sv...")

	r := router.SetupRouter()

	fmt.Println("verificando rotas registradas")

	port := config.GetEnvVar("PORT")
	if port == "" {
		port = ":8080"
	}

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
