package platform

import (
	"fmt"

	email "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	host, username, password, from string
	port                           int
}

func NewMailer(host string, port int, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	client := email.NewSMTPClient()
	client.Host = m.host
	client.Port = m.port
	client.Username = m.username
	client.Password = m.password

	if m.port == 465 {
		client.Encryption = email.EncryptionSSL
	} else {
		client.Encryption = email.EncryptionSTARTTLS
	}

	smtpClient, err := client.Connect()
	if err != nil {
		return err
	}
	defer smtpClient.Close()

	msg := email.NewMSG()
	msg.SetFrom(m.from)
	msg.AddTo(to)
	msg.SetSubject(subject)
	msg.SetBody(email.TextHTML, body)

	if err := msg.Send(smtpClient); err != nil {
		return fmt.Errorf("smtp connect: %w", err)
	}

	return nil
}
