package mailer

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

type Mailer interface {
	SendResetPassword(to string, token string) error
}

type smtpMailer struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

func NewMailer(host string, port int, user string, password string, from string) Mailer {
	return &smtpMailer{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

func (m *smtpMailer) SendResetPassword(to string, token string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Восстановление пароля — GoDelivery")
	msg.SetBody("text/html", fmt.Sprintf(`
		<h2>Восстановление пароля</h2>
		<p>Вы запросили восстановление пароля.</p>
		<p>Ваш токен для сброса пароля:</p>
		<h3>%s</h3>
		<p>Токен действителен 1 час.</p>
		<p>Если вы не запрашивали сброс пароля — проигнорируйте это письмо.</p>
	`, token))

	dialer := gomail.NewDialer(m.host, m.port, m.user, m.password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(msg)
}
