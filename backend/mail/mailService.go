package mail

import (
	"fmt"
	"mime"
	"net/smtp"

	"github.com/gabsfranca/mensagensAnonimasRH/backend/config"
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
	err := config.LoadEnvVars()
	if err != nil {
		return nil, err
	}

	return &EmailServiceSMTP{
		SMTPHost:     config.GetEnvVar("SMTP_HOST"),
		SMTPPort:     config.GetEnvVar("SMTP_PORT"),
		SMTPUsername: config.GetEnvVar("SMTP_USERNAME"),
		SMTPPassword: config.GetEnvVar("SMTP_PASSWORD"),
		Recipient:    config.GetEnvVar("RECIPIENT_EMAIL"),
	}, nil
}

//aqui é tipo os métodos

func (s *EmailServiceSMTP) SendMail(subject, body string) error {
	auth := smtp.PlainAuth("", s.SMTPUsername, s.SMTPPassword, s.SMTPHost)

	from := s.SMTPUsername
	to := []string{s.Recipient}

	//codificando pra utf8 pra poder ter ç

	subjectEncoded := "Assunto: " + mime.QEncoding.Encode("utf-8", subject)
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
