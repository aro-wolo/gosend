package gosend

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

// Environment defines the type for execution mode
type Environment int

const (
	Debug Environment = iota // Debug mode (logs emails instead of sending)
	Test                     // Send mail but using local settings
	Live                     // Live mode (actually sends emails)
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

// Template represents an email template
type Template struct {
	Subject string
	Body    string
}

// ParseTemplate parses an email template from a file
func ParseTemplate(filePath string) (*Template, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %v", err)
	}

	tmpl := &Template{
		Body: string(content),
	}
	return tmpl, nil
}

// RenderTemplate renders the template with the provided data
func (t *Template) RenderTemplate(data interface{}) (string, error) {
	tmpl, err := template.New("email").Parse(t.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var renderedBody string
	writer := &strings.Builder{}
	err = tmpl.Execute(writer, data)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}
	renderedBody = writer.String()

	return renderedBody, nil
}

// Now sends an email using the provided SMTP configuration
func Now(config SMTPConfig, r Recipients, subject, msg string, from ...string) error {
	// Validate SMTP configuration
	if config.Username == "" || config.Password == "" || config.Server == "" {
		return errors.New("SMTP configuration is missing required fields")
	}

	// Set default SMTP port if not provided
	if config.Port == 0 {
		config.Port = 587
	}

	// Validate recipients
	if len(r.To) == 0 {
		return errors.New("no primary recipient specified")
	}

	// Determine sender email
	finalFrom := config.From
	if len(from) > 0 && from[0] != "" {
		finalFrom = from[0]
	} else if finalFrom == "" {
		finalFrom = config.Username
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", finalFrom)
	m.SetHeader("To", r.To...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	// Add CC and BCC if present
	if len(r.Cc) > 0 {
		m.SetHeader("Cc", r.Cc...)
	}
	if len(r.Bcc) > 0 {
		m.SetHeader("Bcc", r.Bcc...)
	}

	// Configure SMTP Dialer
	d := gomail.NewDialer(config.Server, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{
		ServerName:         config.Server,
		InsecureSkipVerify: config.Mode == Test || config.Mode == Debug,
	}

	// Check environment mode
	if config.Mode == Live || config.Mode == Test {
		// Send email in Live mode
		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	} else {
		// Log email details in Debug mode
		log.Printf("'🛠 [DEBUG MODE] Email not sent. Details:\nFrom: %s\nTo: %+v\nCC: %+v\nBCC: %+v\nSubject: %s\n", finalFrom, r.To, r.Cc, r.Bcc, subject)
	}

	return nil
}
