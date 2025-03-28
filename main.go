package gosend

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log"
	"strings"

	"gopkg.in/gomail.v2"
)

// Environment defines execution modes
type Environment int

const (
	Debug Environment = iota // Logs emails instead of sending
	Test                     // Sends mail using local settings
	Live                     // Actually sends emails
)

// SMTPConfig stores SMTP settings
type SMTPConfig struct {
	Username string
	Password string
	Server   string
	Port     int
	Mode     Environment // Debug or Live
	From     string      // Optional sender email
}

// Recipients represents the recipients of the email
type Recipients struct {
	To  []string `json:"to"`
	Cc  []string `json:"cc,omitempty"`
	Bcc []string `json:"bcc,omitempty"`
}

// TemplateManager handles multiple templates
type TemplateManager struct {
	templates *template.Template
}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{}
}

// ParseTemplate loads multiple template files
func (tm *TemplateManager) ParseTemplate(filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("no template files provided")
	}

	tmpl, err := template.ParseFiles(filePaths...)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %v", err)
	}

	tm.templates = tmpl
	return nil
}

// RenderTemplate applies the data struct to all templates
func (tm *TemplateManager) RenderTemplate(data interface{}) (string, error) {
	if tm.templates == nil {
		return "", errors.New("templates not loaded")
	}

	var renderedBody strings.Builder
	err := tm.templates.ExecuteTemplate(&renderedBody, "body.html", data)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}

	return renderedBody.String(), nil
}

// Now sends an email using SMTP
func Now(config SMTPConfig, r Recipients, subject, msg string, from ...string) error {
	if config.Username == "" || config.Password == "" || config.Server == "" {
		return errors.New("SMTP configuration is missing required fields")
	}

	if config.Port == 0 {
		config.Port = 587
	}

	if len(r.To) == 0 {
		return errors.New("no primary recipient specified")
	}

	finalFrom := config.From
	if len(from) > 0 && from[0] != "" {
		finalFrom = from[0]
	} else if finalFrom == "" {
		finalFrom = config.Username
	}

	m := gomail.NewMessage()
	m.SetHeader("From", finalFrom)
	m.SetHeader("To", r.To...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	if len(r.Cc) > 0 {
		m.SetHeader("Cc", r.Cc...)
	}
	if len(r.Bcc) > 0 {
		m.SetHeader("Bcc", r.Bcc...)
	}

	d := gomail.NewDialer(config.Server, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{
		ServerName:         config.Server,
		InsecureSkipVerify: config.Mode == Test || config.Mode == Debug,
	}

	if config.Mode == Live || config.Mode == Test {
		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	} else {
		log.Printf("'🛠 [DEBUG MODE] Email not sent. Details:\nFrom: %s\nTo: %+v\nCC: %+v\nBCC: %+v\nSubject: %s\n", finalFrom, r.To, r.Cc, r.Bcc, subject)
	}

	return nil
}
