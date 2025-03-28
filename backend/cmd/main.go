package main

import (
	"fmt"

	"github.com/gabsfranca/mensagensAnonimasRH/backend/config"
	"github.com/gabsfranca/mensagensAnonimasRH/backend/handler"
	"github.com/gabsfranca/mensagensAnonimasRH/backend/mail"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.LoadEnvVars()
	if err != nil {
		fmt.Println("Erro .env: ", err)
		return
	}

	emailService, err := mail.NewEmailServiceSMTP()
	if err != nil {
		fmt.Println("Erro estanciando o servico de email: ", err)
		return
	}

	r := gin.Default()

	templatePath := "C:\\Users\\Gabriel Menegasso\\Desktop\\gabriel\\programas\\trabalhos\\mensagensAnonimas\\teste2\\frontend\\index.html"

	r.LoadHTMLFiles(templatePath)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Mensagem anonima",
		})
	})

	//middleware

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/send-anonymous-message", func(c *gin.Context) {
		handler.HandleAnonymousMessage(c, emailService)
	})

	port := config.GetEnvVar("PORT")

	if port == "" {
		port = "8080"
	}
	fmt.Printf(port)
	r.Run(":" + port)

}
