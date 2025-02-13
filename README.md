# GoSend

GoSend is a lightweight and flexible Go package for sending emails using SMTP. It supports multiple modes of operation (Debug, Test, and Live), optional email templates for dynamic content generation, and easy configuration for sending emails to multiple recipients, including CC and BCC.

## Features

- **Multiple Environment Modes**:

  - **Debug Mode**: Logs email details instead of sending the email.
  - **Test Mode**: Sends emails but skips TLS verification (useful for local testing).
  - **Live Mode**: Sends emails using the provided SMTP settings.

- **Optional Template System**: Supports dynamic email content generation using Go's `html/template` package.

- **Flexible Configuration**: Allows configuration of SMTP settings, including username, password, server, port, and mode.

- **Customizable Sender**: Supports setting a custom sender email address, with a fallback to the SMTP username if not provided.

- **Recipient Management**: Supports sending emails to multiple recipients, including CC and BCC fields.

## Installation

To use GoSend in your Go project, install it using `go get`:

```bash
go get github.com/aro-wolo/gosend
```

## Usage

### Basic Example (Without Templates)

Here's a simple example of how to use GoSend to send an email without templates:

```go
package main

import (
	"log"

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

	// Send email
	err := gosend.Now(config, recipients, "Plain Subject", "<h1>Hello, World!</h1>")
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
```

### Using Templates

GoSend supports dynamic email content generation using templates. Here's how to use it:

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

### Configuration

The `SMTPConfig` struct is used to configure the SMTP settings:

```go
type SMTPConfig struct {
	Username string      // SMTP username (usually your email address)
	Password string      // SMTP password
	Server   string      // SMTP server address
	Port     int         // SMTP port (defaults to 587 if not provided)
	Mode     Environment // Debug, Test, or Live mode
	From     string      // Optional sender email address
}
```

### Recipients

The `Recipients` struct is used to specify the email recipients:

```go
type Recipients struct {
	To  []string `json:"to"`           // Primary recipients
	Cc  []string `json:"cc,omitempty"` // CC recipients (optional)
	Bcc []string `json:"bcc,omitempty"` // BCC recipients (optional)
}
```

### Template System

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

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/aro-wolo/gosend). If you'd like to contribute code, feel free to fork the repository and submit a pull request.

## License

GoSend is licensed under the MIT License. See the [LICENSE](https://github.com/aro-wolo/gosend/blob/main/LICENSE) file for more details.

## Acknowledgments

- This package uses the [gomail](https://github.com/go-gomail/gomail) library for sending emails.
- Special thanks to all contributors and users of the package.

---

Feel free to explore the [GitHub repository](https://github.com/aro-wolo/gosend) for more details and to contribute to the project!
