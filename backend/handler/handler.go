package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gabsfranca/mensagensAnonimasRH/mail"
	"github.com/gin-gonic/gin"
)

type AnonymousMessage struct {
	Content string `json:"content" binding:"required"`
}

func sanitizeMessage(message string) string {
	message = strings.TrimSpace(message)
	if len(message) > 1000 {
		message = message[:1000]
	}
	return message
}

func HandleAnonymousMessage(c *gin.Context, emailService mail.EmailService) {
	var msg AnonymousMessage

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "mensagem invalida",
		})
		return
	}

	msg.Content = sanitizeMessage(msg.Content)

	if msg.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "a mensagem nao pode estar vazia",
		})
	}

	err := emailService.SendMail("", msg.Content)
	if err != nil {
		log.Println("Erro no envio do email: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "falha ao enviar email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mensagem enviada com sucesso!!!",
	})
}
