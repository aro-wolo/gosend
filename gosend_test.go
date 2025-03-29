package gosend

import (
	"html/template"
	"strings"
	"testing"
)

// Mock SMTP Config
var testSMTPConfig = SMTPConfig{
	Username: "test@example.com",
	Password: "password",
	Server:   "smtp.example.com",
	Port:     587,
	Mode:     Debug,
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
	config := SMTPConfig{}
	err := Now(config, testRecipients, "Test Subject", "Test Message")
	if err == nil || !strings.Contains(err.Error(), "SMTP configuration is missing required fields") {
		t.Errorf("Expected error for missing SMTP configuration, got: %v", err)
	}
}

// Test Now function for missing recipient
func TestNow_MissingRecipient(t *testing.T) {
	recipients := Recipients{}
	err := Now(testSMTPConfig, recipients, "Test Subject", "Test Message")
	if err == nil || !strings.Contains(err.Error(), "no primary recipient specified") {
		t.Errorf("Expected error for missing recipients, got: %v", err)
	}
}

// Test Now function in Debug mode
func TestNow_DebugMode(t *testing.T) {
	config := testSMTPConfig
	config.Mode = Debug

	err := Now(config, testRecipients, "Test Subject", "Test Message")
	if err != nil {
		t.Errorf("Expected no error in debug mode, but got: %v", err)
	}
}

// Test Now function with a custom sender
func TestNow_CustomSender(t *testing.T) {
	customSender := "custom@example.com"
	err := Now(testSMTPConfig, testRecipients, "Test Subject", "Test Message", customSender)
	if err != nil {
		t.Errorf("Expected no error with custom sender, got: %v", err)
	}
}

// Test ParseTemplate function for file reading
/* func TestParseTemplate_Success(t *testing.T) {
	tm := NewTemplateManager()

	// Create a temp file
	tmpFile, err := os.CreateTemp("", "template_*.html")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write a template with an explicit name "body.html"
	sampleTemplate := `{{define "body.html"}}<h1>Hello {{.Name}}</h1>{{end}}`
	tmpFile.WriteString(sampleTemplate)
	tmpFile.Close()

	// Parse template
	err = tm.ParseTemplate(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseTemplate failed: %v", err)
	}

	// Render the template
	data := struct {
		Name string
	}{Name: "John"}

	result, err := tm.RenderTemplate(data)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "<h1>Hello John</h1>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
} */

// Test ParseTemplate function for file not found
func TestParseTemplate_FileNotFound(t *testing.T) {
	tm := NewTemplateManager()

	err := tm.ParseTemplate("nonexistent.html")
	if err == nil || !strings.Contains(err.Error(), "failed to parse templates") {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

// Test RenderTemplate function
/* func TestRenderTemplate_Success(t *testing.T) {
	tm := NewTemplateManager()
	tm.templates = template.Must(template.New("test").Parse("Hello {{.Name}}"))

	data := struct {
		Name string
	}{Name: "John"}

	result, err := tm.RenderTemplate(data)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Hello John"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
} */

// Test RenderTemplate with missing templates
func TestRenderTemplate_NoTemplate(t *testing.T) {
	tm := NewTemplateManager()

	_, err := tm.RenderTemplate(map[string]string{"Name": "John"})
	if err == nil || !strings.Contains(err.Error(), "templates not loaded") {
		t.Errorf("Expected templates not loaded error, got: %v", err)
	}
}

// Test RenderTemplate with invalid template syntax
func TestRenderTemplate_InvalidSyntax(t *testing.T) {
	tm := NewTemplateManager()
	tm.templates = template.Must(template.New("test").Parse("Hello {{.Invalid}}"))

	data := struct {
		Name string
	}{Name: "John"}

	_, err := tm.RenderTemplate(data)
	if err == nil {
		t.Errorf("Expected template parsing error, got nil")
	}
}

// Test Now function with missing subject and message
/* func TestNow_MissingSubjectAndMessage(t *testing.T) {
	err := Now(testSMTPConfig, testRecipients, "", "")
	if err != nil {
		t.Errorf("Expected no error for empty subject and message, got: %v", err)
	}
}
*/
