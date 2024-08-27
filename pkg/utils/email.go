package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, html string) error {

	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")
	smtpUser := os.Getenv("EMAIL_USERNAME")
	smtpPass := os.Getenv("EMAIL_PASSWORD")

	// Create message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		html + "\r\n")

	// Create a TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, // skip SSL verification
		ServerName:         smtpHost,
	}

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Create an SMTP client from the connection
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Authenticate with SMTP server
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %v", err)
	}

	// Set sender and recipient
	if err = client.Mail(smtpUser); err != nil {
		return fmt.Errorf("failed to set sender email: %v", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient email: %v", err)
	}

	// Send email data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send email data: %v", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email message: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close email writer: %v", err)
	}

	// Close SMTP connection
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to close SMTP connection: %v", err)
	}

	return nil
}
