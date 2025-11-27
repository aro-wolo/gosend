package main

import (
	"log"

	"github.com/aro-wolo/gosend/v2" // update path to your module
)

func main() {
	config := gosend.SMTPConfig{
		Username: "your-email@example.com",
		Password: "your-email-password",
		Server:   "smtp.example.com",
		Port:     587,
		Mode:     gosend.Debug, // Debug logs without sending
		From:     "no-reply@example.com",
	}

	recipients := gosend.Recipients{
		To:  []string{"recipient1@example.com"},
		Cc:  []string{"cc@example.com"},  // optional
		Bcc: []string{"bcc@example.com"}, // optional
	}

	subject := "Hello from GoSend"
	message := `
	<h1>Hi there!</h1>
	<p>This is a test email sent using GoSend's Now() function without templates.</p>
	`

	err := gosend.Now(config, recipients, subject, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
