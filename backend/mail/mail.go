package mail

//abstraindo os meus services de email (interface abstrata)
type EmailService interface {
	SendMail(subject, body string) error
}
