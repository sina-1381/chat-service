package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type smtpServer struct {
	host string
	port string
}

type SetEmail struct {
	To []string
	Subject string
	Massage string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func (s *SetEmail)SendEmail() {
	from := "sina.darestany@gmail.com"
	password := "siNadarestaNy1381"
	to := s.To
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(
		"Subject: "+s.Subject+"!\n" +
			mime + "\n"+
			s.Massage)
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

func TemplateEmailSender(url , token , input string )  {
	var tpl bytes.Buffer
	template:=template.Must(template.ParseFiles("templates/emailResetPassword.tmpl"))
	params := map[string]interface{}{
		"action_url": "http://localhost:8080+"+url+"/"+token,
	}
	if err := template.Execute(&tpl, params); err != nil {
		panic(err)
	}
	email := SetEmail{
		To:      []string{input},
		Subject: "Reset Password(ginAuth)",
		Massage: tpl.String(),
	}
	email.SendEmail()
}