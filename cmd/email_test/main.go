package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

func main() {
	smtp_account := requireNonNil("SMTP_ACCOUNT")
	smtp_password := requireNonNil("SMTP_PASSWORD")
	smtp_host := requireNonNil("SMTP_HOST")
	smtp_port_string := requireNonNil("SMTP_PORT")
	smtp_from := requireNonNil("SMTP_FROM")
	smtp_port, err := strconv.Atoi(smtp_port_string)
	if err != nil {
		panic("SMTP_PORT is not a valid port number")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtp_from)
	m.SetHeader("To", os.Getenv("TO_EMAIL"))
	m.SetHeader("This is a test from DigitalVenue")
	m.SetBody("text/plain", "This is a test email from DigitalVenue (digital-venue.net)")

	d := gomail.NewDialer(smtp_host, smtp_port, smtp_account, smtp_password)
	err = d.DialAndSend(m)
	if err != nil {
		panic("failed to send email: " + err.Error())
	}

	fmt.Println("Email sent")
}

func requireNonNil(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("missing required environment variable: " + key)
	}
	return val
}
