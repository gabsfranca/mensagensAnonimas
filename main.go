package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func sendAnonymousEmail(message string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("erro ao carregar o arquivo .env: %v", err)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	recipientEmail := os.Getenv("RECIPIENT_EMAIL")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || recipientEmail == "" {
		return fmt.Errorf("configurações de email incompletas")
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	from := smtpUsername
	to := []string{recipientEmail}

	subject := "Assunto: " + mime.QEncoding.Encode("utf-8", "Mensagem Anônima")
	body := fmt.Sprintf("\n\n%s", message)

	emailMessage := []byte("To: " + recipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, emailMessage)
	if err != nil {
		return fmt.Errorf("erro ao enviar e-mail: %v", err)
	}

	return nil
}

func handleAnonymousMessage(c *gin.Context) {
	var msg AnonymousMessage

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Mensagem inválida",
		})
		return
	}

	msg.Content = sanitizeMessage(msg.Content)

	if msg.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "A mensagem não pode estar vazia",
		})
		return
	}

	err := sendAnonymousEmail(msg.Content)
	if err != nil {
		log.Println("Erro no envio de email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao enviar email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mensagem enviada com sucesso",
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

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

	r.Static("/", "./frontend")

	r.POST("/send-anonymous-message", handleAnonymousMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor iniciando na porta %s\n", port)
	r.Run(":" + port)
}
