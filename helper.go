package gosend

import (
	"fmt"
)

// SendMail is a convenience wrapper that loads header/body/footer templates,
// renders them, and sends a fully formatted email.
func SendMail(
	config SMTPConfig,
	recipients Recipients,
	subject string,
	templateName string,
	data map[string]any,
	templateBasePath ...string, // optional override
) error {

	// Default template folder
	base := "templates"
	if len(templateBasePath) > 0 && templateBasePath[0] != "" {
		base = templateBasePath[0]
	}

	tm := NewTemplateManager()

	err := tm.ParseTemplate(
		fmt.Sprintf("%s/header.html", base),
		fmt.Sprintf("%s/%s.html", base, templateName),
		fmt.Sprintf("%s/footer.html", base),
	)

	if err != nil {
		return fmt.Errorf("failed parsing templates: %w", err)
	}

	body, err := tm.RenderTemplate(data)
	if err != nil {
		return fmt.Errorf("failed rendering template: %w", err)
	}

	// send
	if err := Now(config, recipients, subject, body); err != nil {
		return fmt.Errorf("failed sending email: %w", err)
	}

	return nil
}
