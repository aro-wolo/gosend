package gosend

import (
	"os"
	"strings"
	"testing"
)

// Mock SMTP Config
var testSMTPConfig = SMTPConfig{
	Username: "test@example.com",
	Password: "password",
	Server:   "smtp.example.com",
	Port:     587,
	Mode:     Live,
	From:     "sender@example.com",
}

// Mock Recipients
var testRecipients = Recipients{
	To:  []string{"recipient@example.com"},
	Cc:  []string{"cc@example.com"},
	Bcc: []string{"bcc@example.com"},
}

// Test Now function for missing SMTP config
func TestNow_MissingSMTPConfig(t *testing.T) {
	config := SMTPConfig{} // Missing all required fields
	err := Now(config, testRecipients, "Test Subject", "Test Message")
	if err == nil || err.Error() != "SMTP configuration is missing required fields" {
		t.Errorf("Expected error 'SMTP configuration is missing required fields', got: %v", err)
	}
}

// Test Now function for missing recipient
func TestNow_MissingRecipient(t *testing.T) {
	config := testSMTPConfig
	recipients := Recipients{}
	err := Now(config, recipients, "Test Subject", "Test Message")
	if err == nil || err.Error() != "no primary recipient specified" {
		t.Errorf("Expected error 'no primary recipient specified', got: %v", err)
	}
}

// Test Now function in Debug mode (should log and not send)
func TestNow_DebugMode(t *testing.T) {
	config := testSMTPConfig
	config.Mode = Debug
	err := Now(config, testRecipients, "Test Subject", "Test Message")
	if err != nil {
		t.Errorf("Expected no error in debug mode, but got: %v", err)
	}
}

/*
// Test Now function in Test mode (should not verify TLS)
func TestNow_TestMode(t *testing.T) {
	config := testSMTPConfig
	config.Mode = Test
	err := Now(config, testRecipients, "Test Subject", "Test Message")
	if err != nil {
		t.Errorf("Expected no error in test mode, but got: %v", err)
	}

} */

/*
// Test Now function in Live mode (should attempt to send)
func TestNow_LiveMode(t *testing.T) {
	config := testSMTPConfig
	config.Mode = Live
	err := Now(config, testRecipients, "Test Subject", "Test Message")
	if err != nil {
		t.Errorf("Expected no error in live mode, but got: %v", err)
	}
}
*/
// Test ParseTemplate function for file reading
func TestParseTemplate_Success(t *testing.T) {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "template_*.html")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write sample template
	sampleTemplate := "<h1>Hello {{.Name}}</h1>"
	tmpFile.WriteString(sampleTemplate)
	tmpFile.Close()

	tmpl, err := ParseTemplate(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	if tmpl.Body != sampleTemplate {
		t.Errorf("Expected template body to be '%s', got '%s'", sampleTemplate, tmpl.Body)
	}
}

// Test ParseTemplate function for file not found
func TestParseTemplate_FileNotFound(t *testing.T) {
	_, err := ParseTemplate("nonexistent.html")
	if err == nil || !strings.Contains(err.Error(), "failed to read template file") {
		t.Errorf("Expected 'failed to read template file' error, got: %v", err)
	}
}

// Test RenderTemplate function
func TestRenderTemplate_Success(t *testing.T) {
	tmpl := &Template{
		Body: "Hello {{.Name}}",
	}

	data := struct {
		Name string
	}{Name: "John"}

	result, err := tmpl.RenderTemplate(data)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Hello John"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Test RenderTemplate with invalid template syntax
func TestRenderTemplate_InvalidSyntax(t *testing.T) {
	tmpl := &Template{
		Body: "Hello {{.x}}",
	}

	data := struct {
		Name string
	}{Name: "John"}

	_, err := tmpl.RenderTemplate(data)
	if err == nil {
		t.Errorf("Expected template parsing error, got nil")
	}
}
