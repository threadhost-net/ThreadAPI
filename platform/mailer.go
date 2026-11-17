package platform

import (
	"fmt"

	email "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	Host, Username, Password, From string
	Port                           int
}

func (m *Mailer) Send(to, subject, body string) error {
	client := email.NewSMTPClient()
	client.Host = m.Host
	client.Port = m.Port
	client.Username = m.Username
	client.Password = m.Password

	if m.Port == 465 {
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
	msg.SetFrom("ThreadHost <" + m.From + ">")
	msg.AddTo(to)
	msg.SetSubject(subject)
	msg.SetBody(email.TextHTML, body)

	if err := msg.Send(smtpClient); err != nil {
		return fmt.Errorf("smtp connect: %w", err)
	}

	return nil
}
