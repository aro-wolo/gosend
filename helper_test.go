package gosend

import (
	"os"
	"testing"
)

func TestSendMailDebug(t *testing.T) {
	// Create a temporary directory for templates
	tmpDir := "examples/templates"
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	// Clean up after test
	defer os.RemoveAll("examples/templates")

	// Create template files
	files := map[string]string{
		"header.html":  "<h1>Welcome to GoSend</h1>",
		"welcome.html": "<p>Hello {{.Name}}</p>",
		"footer.html":  "<hr>Thank you",
	}

	for name, content := range files {
		if err = os.WriteFile(tmpDir+"/"+name, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template file %s: %v", name, err)
		}
	}

	// SMTP configuration in Debug mode
	config := SMTPConfig{
		Username: "dummy",
		Password: "dummy",
		Server:   "smtp.test",
		Mode:     Debug,
	}

	// Recipient list
	recipients := Recipients{
		To: []string{"user@test.com"},
	}

	err = SendMail(config, recipients, "Test Subject", "welcome", map[string]any{
		"Name": "Tester Guy",
	}, tmpDir)

	if err != nil {
		t.Errorf("SendMail failed: %v", err)
	}
}
