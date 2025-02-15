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

# GoSend Email Package

This Go package provides a simple and flexible way to send emails using SMTP. It supports multiple environments (Debug, Test, and Live), allowing you to log emails instead of sending them during development or testing. The package is built on top of the popular `gomail.v2` library.

---

## Features

- **Multiple Environments**: Supports Debug, Test, and Live modes.
- **Flexible Configuration**: Allows custom SMTP settings, including username, password, server, and port.
- **Recipient Management**: Supports To, Cc, and Bcc recipients.
- **HTML Emails**: Sends emails with HTML content.
- **TLS Support**: Configurable TLS settings for secure email delivery.
- **Debug Mode**: Logs email details instead of sending them in Debug mode.
- **Template Support**: Optional template system for dynamic email content generation.

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

1. **Create a Template File**:
   Save the following content in a file named `notification_template.html`:

   ```html
   <!DOCTYPE html>
   <html>
   	<head>
   		<title>Notification</title>
   	</head>
   	<body>
   		<h1>Hello, {{.Name}}!</h1>
   		<p>This is a notification regarding your recent activity.</p>
   		<p>Details:</p>
   		<ul>
   			<li>Activity: {{.Activity}}</li>
   			<li>Date: {{.Date}}</li>
   		</ul>
   		<p>If you have any questions, feel free to contact us.</p>
   		<p>Best regards,</p>
   		<p>The Team</p>
   	</body>
   </html>
   ```

2. **Send Email Using the Template**:

   ```go
   package main

   import (
   	"log"
   	"time"

   	"github.com/aro-wolo/gosend"
   )

   func main() {
   	// Define SMTP configuration
   	config := gosend.SMTPConfig{
   		Username: "your-email@example.com",
   		Password: "your-email-password",
   		Server:   "smtp.example.com",
   		Port:     587,
   		Mode:     gosend.Live,
   		From:     "your-email@example.com",
   	}

   	// Define recipients
   	recipients := gosend.Recipients{
   		To: []string{"recipient@example.com"},
   	}

   	// Parse the template
   	template, err := gosend.ParseTemplate("notification_template.html")
   	if err != nil {
   		log.Fatalf("Failed to parse template: %v", err)
   	}

   	// Define template data
   	data := struct {
   		Name     string
   		Activity string
   		Date     string
   	}{
   		Name:     "John Doe",
   		Activity: "Account Login",
   		Date:     time.Now().Format("2006-01-02 15:04:05"),
   	}

   	// Render the template
   	renderedBody, err := template.RenderTemplate(data)
   	if err != nil {
   		log.Fatalf("Failed to render template: %v", err)
   	}

   	// Send email
   	err = gosend.Now(config, recipients, "Notification", renderedBody)
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
| ---------- | --------------------------------------------------- | ------------- |
| `Username` | SMTP username (usually your email address).         | Required      |
| `Password` | SMTP password.                                      | Required      |
| `Server`   | SMTP server address (e.g., `smtp.gmail.com`).       | Required      |
| `Port`     | SMTP port (e.g., `587` for TLS).                    | `587`         |
| `Mode`     | Execution mode: `Debug`, `Test`, or `Live`.         | `Live`        |
| `From`     | Sender email address. If empty, `Username` is used. | Optional      |

### `Recipients`

| Field | Description                              | Example                             |
| ----- | ---------------------------------------- | ----------------------------------- |
| `To`  | Primary recipients (required).           | `[]string{"recipient@example.com"}` |
| `Cc`  | Carbon copy recipients (optional).       | `[]string{"cc@example.com"}`        |
| `Bcc` | Blind carbon copy recipients (optional). | `[]string{"bcc@example.com"}`       |

---

## Modes

### 1. **Debug Mode**

- Logs email details instead of sending them.
- Useful for development and testing.
- Example:
  ```go
  config.Mode = gosend.Debug
  ```

### 2. **Test Mode**

- Sends emails but skips TLS certificate verification.
- Useful for testing with self-signed certificates.
- Example:
  ```go
  config.Mode = gosend.Test
  ```

### 3. **Live Mode**

- Sends emails with full TLS verification.
- Use this in production.
- Example:
  ```go
  config.Mode = gosend.Live
  ```

---

## Template System

The `Template` struct and its methods allow for dynamic email content generation:

```go
type Template struct {
	Subject string
	Body    string
}

// ParseTemplate parses an email template from a file
func ParseTemplate(filePath string) (*Template, error)

// RenderTemplate renders the template with the provided data
func (t *Template) RenderTemplate(data interface{}) (string, error)
```

---

## Error Handling

The `Now` function returns an error if:

- SMTP configuration is incomplete.
- No primary recipient is specified.
- Email sending fails.

Example error handling:

```go
err := gosend.Now(config, recipients, subject, message)
if err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
```

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- Built using [gomail.v2](https://github.com/go-gomail/gomail).
- Inspired by the need for a simple and flexible email-sending solution in Go.

---

## Support

For questions or issues, please open an issue on [GitHub](https://github.com/aro-wolo/gosend/issues).

---
