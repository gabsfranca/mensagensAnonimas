package router

import (
	"log"
	"net/http"

	"github.com/gabsfranca/mensagensAnonimasRH/handler"
	"github.com/gabsfranca/mensagensAnonimasRH/mail"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") //usando * por ser em desenvolvimento
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	emailService, err := mail.NewEmailServiceSMTP()
	if err != nil {
		log.Println("erro estancioando emailService: ", err)
	}

	r.POST("/send-anonymous-message", func(c *gin.Context) {
		handler.HandleAnonymousMessage(c, emailService)
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Mensagem anônima",
		})
	})

	return r

}
