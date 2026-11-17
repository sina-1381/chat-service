package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type smtpServer struct {
	host string
	port string
}

type SetEmail struct { // email content structure
	To      []string
	Subject string
	Massage string
}

func (s *smtpServer) Address() string { // returns smtp server address
	return s.host + ":" + s.port
}

func (s *SetEmail) SendEmail() { // email send function

	// info from .env file
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASS")
	to := s.To // users email address

	// smtp connection
	smtpServer := smtpServer{host: os.Getenv("EMAIL_SERVER_HOST"), port: os.Getenv("EMAIL_SERVER_PORT")}

	// email subject and message content
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(
		"Subject: " + s.Subject + "!\n" +
			mime + "\n" +
			s.Massage)

	// smtp authentication with email password details
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	// sending email to smtp provider
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

func TemplateEmailSender(url, token, input string) { // template email send function
	var tpl bytes.Buffer

	// loading the temp file
	template := template.Must(template.ParseFiles("templates/emailResetPassword.tmpl"))

	// verify reset password link provided in reset password email with random string token
	params := map[string]interface{}{
		"action_url": "http://localhost:8080+" + url + "/" + token,
	}
	if err := template.Execute(&tpl, params); err != nil {
		panic(err)
	}

	// sending email by calling the send email function that handels smtp
	email := SetEmail{
		To:      []string{input},
		Subject: "Reset Password(ginAuth)",
		Massage: tpl.String(),
	}
	email.SendEmail()
}
