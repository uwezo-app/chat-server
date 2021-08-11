package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
	"gorm.io/gorm"
)

func SendEmail(dbase *gorm.DB, email string, HTMLtemp string, body interface{}, writer *bytes.Buffer) error {
	// var user *db.Psychologist

	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")

	to := []string{
		email,
	}

	m := mail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(from, "Uwezo Team")},
		"To":      to,
		"Subject": {"Reset Password"},
	})

	t, _ := template.ParseFiles(HTMLtemp)

	err := t.Execute(writer, body)
	if err != nil {
		return err
	}

	m.SetBody("text/html", writer.String())

	p, _ := strconv.Atoi(port)
	d := mail.NewDialer(host, p, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = d.DialAndSend(m); err != nil {
		return errors.New("could not send you an email")
	}

	return nil
}
