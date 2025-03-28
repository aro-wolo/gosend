<!-- [![Go Reference](https://pkg.go.dev/badge/github.com/aro-wolo/gosend.svg)](https://pkg.go.dev/github.com/aro-wolo/gosend) -->

![GitHub Repo stars](https://img.shields.io/github/stars/aro-wolo/gosend?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/aro-wolo/gosend)
![GitHub issues](https://img.shields.io/github/issues/aro-wolo/gosend)
![GitHub contributors](https://img.shields.io/github/contributors/aro-wolo/gosend)
[![Go Report Card](https://goreportcard.com/badge/github.com/aro-wolo/gosend)](https://goreportcard.com/report/github.com/aro-wolo/gosend)
[![Build Status](https://github.com/aro-wolo/gosend/actions/workflows/go.yml/badge.svg)](https://github.com/aro-wolo/gosend/actions)
[![codecov](https://codecov.io/gh/aro-wolo/gosend/branch/main/graph/badge.svg)](https://codecov.io/gh/aro-wolo/gosend)
![Go Module](https://img.shields.io/github/go-mod/go-version/aro-wolo/gosend)
[![License](https://img.shields.io/github/license/aro-wolo/gosend.svg)](https://github.com/aro-wolo/gosend/blob/main/LICENSE)
[![wakatime](https://wakatime.com/badge/user/c78c31fe-9c97-4b21-b865-91bc182f2d42.svg)](https://wakatime.com/@c78c31fe-9c97-4b21-b865-91bc182f2d42)

# GoSend Email Package

This Go package provides a simple and flexible way to send emails using SMTP. It supports multiple environments (Debug, Test, and Live), allowing you to log emails instead of sending them during development or testing. The package is built on top of the popular `gomail.v2` library and includes a powerful template management system for dynamic email content generation.

---

## Features

- **Multiple Environments**: Supports Debug, Test, and Live modes.
- **Flexible Configuration**: Allows custom SMTP settings, including username, password, server, and port.
- **Recipient Management**: Supports To, Cc, and Bcc recipients.
- **HTML Emails**: Sends emails with HTML content.
- **TLS Support**: Configurable TLS settings for secure email delivery.
- **Debug Mode**: Logs email details instead of sending them in Debug mode.
- **Template System**: Supports multi-part templates (header, body, footer) for dynamic email content.

---

## Installation

To use this package in your Go project, run:

```bash
go get github.com/aro-wolo/gosend
```

---

## Usage

### Import the Package

```go
import (
	"github.com/aro-wolo/gosend"
)
```

### Example: Sending an Email Without Templates

```go
package main

import (
	"fmt"
	"log"

	"github.com/aro-wolo/gosend"
)

func main() {
	// SMTP Configuration
	config := gosend.SMTPConfig{
		Username: "your-email@example.com",
		Password: "your-email-password",
		Server:   "smtp.example.com",
		Port:     587,
		Mode:     gosend.Live, // Use Debug, Test, or Live
		From:     "no-reply@example.com",
	}

	// Recipients
	recipients := gosend.Recipients{
		To:  []string{"recipient1@example.com", "recipient2@example.com"},
		Cc:  []string{"cc@example.com"},
		Bcc: []string{"bcc@example.com"},
	}

	// Email Content
	subject := "Test Email"
	message := `<h1>Hello, World!</h1><p>This is a test email.</p>`

	// Send Email
	err := gosend.Now(config, recipients, subject, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
```

### Example: Sending an Email With Templates

1. **Create Template Files**:

   Save the following content as separate files:

   **header.html**
   ```html
   <h1>Welcome to Our Service</h1>
   ```

   **body.html**
   ```html
   <p>Hello {{.Name}},</p>
   <p>We are excited to have you on board!</p>
   ```

   **footer.html**
   ```html
   <p>Thank you for choosing us! <br> Contact us at support@example.com</p>
   ```

2. **Send Email Using the Template System**:

   ```go
   package main

   import (
   	"log"
   	"github.com/aro-wolo/gosend"
   )

   func main() {
   	// Initialize template manager
   	tm := gosend.NewTemplateManager()

   	// Parse templates
   	err := tm.ParseTemplate("header.html", "body.html", "footer.html")
   	if err != nil {
   		log.Fatalf("Failed to parse templates: %v", err)
   	}

   	// Data to inject into the template
   	data := map[string]string{
   		"Name": "John Doe",
   	}

   	// Render template
   	renderedBody, err := tm.RenderTemplate(data)
   	if err != nil {
   		log.Fatalf("Failed to render template: %v", err)
   	}

   	// Define SMTP configuration
   	config := gosend.SMTPConfig{
   		Username: "your-email@example.com",
   		Password: "your-password",
   		Server:   "smtp.example.com",
   		Port:     587,
   		Mode:     gosend.Live,
   		From:     "your-email@example.com",
   	}

   	// Define recipients
   	recipients := gosend.Recipients{
   		To: []string{"recipient@example.com"},
   	}

   	// Send email
   	err = gosend.Now(config, recipients, "Welcome!", renderedBody)
   	if err != nil {
   		log.Fatalf("Failed to send email: %v", err)
   	}

   	log.Println("Email sent successfully!")
   }
   ```

---

## Configuration

### `SMTPConfig`

| Field      | Description                                         | Default Value |
|------------|-----------------------------------------------------|---------------|
| `Username` | SMTP username (usually your email address).         | Required      |
| `Password` | SMTP password.                                      | Required      |
| `Server`   | SMTP server address (e.g., `smtp.gmail.com`).       | Required      |
| `Port`     | SMTP port (e.g., `587` for TLS).                    | `587`         |
| `Mode`     | Execution mode: `Debug`, `Test`, or `Live`.         | `Live`        |
| `From`     | Sender email address. If empty, `Username` is used. | Optional      |

### `Recipients`

| Field | Description                              | Example                             |
|-------|------------------------------------------|-------------------------------------|
| `To`  | Primary recipients (required).           | `[]string{"recipient@example.com"}` |
| `Cc`  | Carbon copy recipients (optional).       | `[]string{"cc@example.com"}`        |
| `Bcc` | Blind carbon copy recipients (optional). | `[]string{"bcc@example.com"}`       |

---

## Modes

### 1. **Debug Mode**

- Logs email details instead of sending them.
- Example:
  ```go
  config.Mode = gosend.Debug
  ```

### 2. **Test Mode**

- Sends emails but skips TLS verification.
- Example:
  ```go
  config.Mode = gosend.Test
  ```

### 3. **Live Mode**

- Sends emails with full TLS verification.
- Example:
  ```go
  config.Mode = gosend.Live
  ```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Support

For questions or issues, please open an issue on [GitHub](https://github.com/aro-wolo/gosend/issues).


