package gosend

import (
	"os"
	"testing"
)

// TestSendMail_Debug tests SendMail in Debug mode
func TestSendMail_Debug(t *testing.T) {
	// Setup temporary template directory
	tmpDir := "tmp/templates"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	defer os.RemoveAll("tmp") // cleanup after test

	// Create dummy templates
	os.WriteFile(tmpDir+"/header.html", []byte("<h1>Header</h1>"), 0644)
	os.WriteFile(tmpDir+"/welcome.html", []byte("<p>Hello {{.Name}}</p>"), 0644)
	os.WriteFile(tmpDir+"/footer.html", []byte("<hr>Footer</hr>"), 0644)

	// Dummy SMTP config in Debug mode
	config := SMTPConfig{
		Username: "dummy",
		Password: "dummy",
		Server:   "smtp.test",
		Mode:     Debug,
	}

	recipients := Recipients{
		To: []string{"user@test.com"},
	}

	// Call SendMail
	err := SendMail(config, recipients, "Test Subject", "welcome", map[string]any{
		"Name": "Tester",
	}, tmpDir)

	if err != nil {
		t.Errorf("SendMail failed: %v", err)
	}
}

// TestSendMail_TemplateNotFound checks that missing templates return error
func TestSendMail_TemplateNotFound(t *testing.T) {
	config := SMTPConfig{
		Username: "dummy",
		Password: "dummy",
		Server:   "smtp.test",
		Mode:     Debug,
	}

	recipients := Recipients{
		To: []string{"user@test.com"},
	}

	// Intentionally missing template
	err := SendMail(config, recipients, "Test", "nonexistent", map[string]any{}, "tmp")

	if err == nil {
		t.Fatal("Expected error for missing template, got nil")
	}
}
