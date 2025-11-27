package main

import (
	"log"

	"github.com/aro-wolo/gosend/v2"
)

func main() {
	config := gosend.SMTPConfig{
		Username: "your-email@example.com",
		Password: "password",
		Server:   "smtp.example.com",
		Port:     587,
		Mode:     gosend.Debug,
	}

	recipients := gosend.Recipients{
		To: []string{"test@example.com"},
	}

	data := map[string]any{
		"Name": "Abiodun",
	}

	err := gosend.SendMail(config, recipients, "Welcome!", "welcome", data, "templates")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Email sent!")
	}
}
