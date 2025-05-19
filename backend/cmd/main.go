package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gabsfranca/mensagensAnonimasRH/config"
	"github.com/gabsfranca/mensagensAnonimasRH/router"
)

func main() {

	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)

	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Erro ao criar arquivo de log: %v", err)
	} else {
		log.SetOutput(f)
		fmt.Println("logs sendo salvos em debug.log")
	}

	fmt.Println("iniciandop sv...")

	err = config.LoadEnvVars()
	if err != nil {
		fmt.Println("Erro .env: ", err)
	}

	r := router.SetupRouter()

	fmt.Println("verificando rotas registradas")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
