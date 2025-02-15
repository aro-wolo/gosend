package gosend_test

import (
	"errors"
	"os"
	"testing"

	"github.com/aro-wolo/gosend"
)

func TestParseTemplate(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_template.html")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	templateContent := "Hello {{.Name}}, welcome!"
	if _, err := tempFile.WriteString(templateContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	tmpl, err := gosend.ParseTemplate(tempFile.Name())
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	if tmpl.Body != templateContent {
		t.Errorf("Expected template body %q, got %q", templateContent, tmpl.Body)
	}
}

func TestRenderTemplate(t *testing.T) {
	tmpl := &gosend.Template{
		Body: "Hello {{.Name}}, welcome!",
	}
	rendered, err := tmpl.RenderTemplate(map[string]string{"Name": "Abiodun"})
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Hello Abiodun, welcome!"
	if rendered != expected {
		t.Errorf("Expected %q, got %q", expected, rendered)
	}
}

func TestNow_DebugMode(t *testing.T) {
	config := gosend.SMTPConfig{
		Username: "test@example.com",
		Password: "password",
		Server:   "smtp.example.com",
		Port:     587,
		Mode:     gosend.Debug,
	}

	recipients := gosend.Recipients{
		To: []string{"receiver@example.com"},
	}

	err := gosend.Now(config, recipients, "Test Subject", "Test Message")
	if err != nil {
		t.Errorf("Expected no error in Debug mode, got: %v", err)
	}
}

func TestNow_MissingConfig(t *testing.T) {
	config := gosend.SMTPConfig{}
	recipients := gosend.Recipients{To: []string{"receiver@example.com"}}

	err := gosend.Now(config, recipients, "Test Subject", "Test Message")
	if err == nil {
		t.Errorf("Expected error due to missing SMTP config, but got nil")
	}
}

func TestNow_MissingRecipient(t *testing.T) {
	config := gosend.SMTPConfig{
		Username: "test@example.com",
		Password: "password",
		Server:   "smtp.example.com",
		Port:     587,
		Mode:     gosend.Live,
	}

	recipients := gosend.Recipients{} // No recipient specified

	err := gosend.Now(config, recipients, "Test Subject", "Test Message")
	if !errors.Is(err, errors.New("no primary recipient specified")) {
		t.Errorf("Expected 'no primary recipient specified' error, got: %v", err)
	}
}
