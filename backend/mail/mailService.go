package mail

import (
	"fmt"
	"mime"
	"net/smtp"
	"os"

	"github.com/gabsfranca/mensagensAnonimasRH/config"
)

// implementação concreta do serviço de email
// aqui é tipo os atributos de uma classe
type EmailServiceSMTP struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	Recipient    string
}

// instanciando o emailservicesmtp
// aqui é tipo o constructor
func NewEmailServiceSMTP() (*EmailServiceSMTP, error) {
	// Tentar carregar variáveis, mas continuar mesmo se falhar
	config.LoadEnvVars()

	// Verificar se as variáveis essenciais estão presentes
	smtpHost := config.GetEnvVar("SMTP_HOST")
	smtpPort := config.GetEnvVar("SMTP_PORT")
	smtpUsername := config.GetEnvVar("SMTP_USERNAME")
	smtpPassword := config.GetEnvVar("SMTP_PASSWORD")
	recipient := config.GetEnvVar("RECIPIENT_EMAIL")

	// Em modo de desenvolvimento/teste, podemos configurar um serviço de e-mail simulado
	if smtpHost == "" || smtpPort == "" || smtpUsername == "" ||
		smtpPassword == "" || recipient == "" {
		fmt.Println("AVISO: Algumas variáveis de email não estão configuradas. Usando modo de simulação.")
		return &EmailServiceSMTP{
			SMTPHost:     "localhost",
			SMTPPort:     "25",
			SMTPUsername: "test@example.com",
			SMTPPassword: "password",
			Recipient:    "recipient@example.com",
		}, nil
	}

	return &EmailServiceSMTP{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		Recipient:    recipient,
	}, nil
}

// aqui é tipo os métodos
func (s *EmailServiceSMTP) SendMail(subject, body string) error {
	// Verificar se estamos em modo simulado (desenvolvimento local)
	if s.SMTPHost == "localhost" && os.Getenv("ENVIRONMENT") != "production" {
		fmt.Println("SIMULAÇÃO DE EMAIL: Enviando para", s.Recipient)
		fmt.Println("Assunto:", subject)
		fmt.Println("Conteúdo:", body)
		return nil
	}

	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)

	from := s.SMTPUsername
	to := []string{s.Recipient}

	//codificando pra utf8 pra poder ter ç
	if subject == "" {
		subject = "Mensagem Anônima"
	}

	subjectEncoded := mime.QEncoding.Encode("utf-8", subject)
	emailMessage := []byte("To: " + s.Recipient + "\r\n" +
		"Subject: " + subjectEncoded + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(s.SMTPHost+":"+s.SMTPPort, auth, from, to, emailMessage)

	if err != nil {
		return fmt.Errorf("erro ao enviar e-mail: %v", err)
	}
	return nil
}
