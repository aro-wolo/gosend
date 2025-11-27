package gosend

import (
	"os"
	"testing"
)

func TestSendMailDebug(t *testing.T) {
	os.Mkdir("examples/templates", 0755)
	os.WriteFile("examples/templates/header.html", []byte("<h1>Welcome to GoSend</h1>"), 0644)
	os.WriteFile("examples/templates/welcome.html", []byte("<p>Hello {{.Name}}</p>"), 0644)
	os.WriteFile("examples/templates/footer.html", []byte("<hr>Thank you"), 0644)

	config := SMTPConfig{
		Username: "dummy",
		Password: "dummy",
		Server:   "smtp.test",
		Mode:     Debug,
	}

	recipients := Recipients{
		To: []string{"user@test.com"},
	}

	err := SendMail(config, recipients, "Test", "welcome", map[string]any{
		"Name": "Tester Guy",
	}, "examples/templates")

	if err != nil {
		t.Errorf("SendMail failed: %v", err)
	}
}
