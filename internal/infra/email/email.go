package email

import (
	conf "github.com/Ablebil/sea-catering-be/config"
	"gopkg.in/gomail.v2"
)

type EmailItf interface {
	SendOTPEmail(to string, otp string) error
}

type Email struct {
	sender   string
	password string
}

func NewEmail(conf *conf.Config) EmailItf {
	return &Email{
		sender:   conf.EmailUser,
		password: conf.EmailPassword,
	}
}

func (e *Email) SendOTPEmail(to string, otp string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", e.sender)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Your OTP Code")
	mail.SetBody("text/plain", "Your OTP code is: "+otp)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, e.sender, e.password)
	return dialer.DialAndSend(mail)
}
