package main

import (
	"fmt"
	"os"

	"github.com/gabsfranca/mensagensAnonimasRH/config"
	"github.com/gabsfranca/mensagensAnonimasRH/router"
)

func main() {
	err := config.LoadEnvVars()
	if err != nil {
		fmt.Println("Erro .env: ", err)
	}

	r := router.SetupRouter()

	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "../frontend/index.html"
	}

	r.Static("/static", "./frontend")
	r.LoadHTMLFiles(templatePath)

	port := config.GetEnvVar("PORT")

	if port == "" {
		port = "8080"
	}
	fmt.Printf("Servidor iniciado na porta %s\n", port)
	r.Run(":" + port)
}
