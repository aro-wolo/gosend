package gosend

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

// Environment defines execution modes
type Environment int

const (
	Debug Environment = iota // Logs emails instead of sending
	Test                     // Sends mail using local settings
	Live                     // Actually sends emails
)

// SMTPConfig stores SMTP settings
type SMTPConfig struct {
	Username string
	Password string
	Server   string
	Port     int
	Mode     Environment // Debug or Live
	From     string      // Optional sender email
}

// Attachment describes a file to send with an email.
type Attachment struct {
	Path      string // File path on disk
	Inline    bool   // False sends a normal attachment; true embeds the file
	ContentID string // Optional inline CID, without angle brackets; defaults to the file's base name
}

// Recipients represents the recipients of the email
type Recipients struct {
	To  []string `json:"to"`
	Cc  []string `json:"cc,omitempty"`
	Bcc []string `json:"bcc,omitempty"`
}

// TemplateManager handles multiple templates
type TemplateManager struct {
	templates *template.Template
	order     []string
}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{}
}

func (tm *TemplateManager) ParseTemplate(filePaths ...string) error {
	if len(filePaths) == 0 {
		return errors.New("no template files provided")
	}

	t := template.New("")
	tm.order = []string{} // reset

	for _, path := range filePaths {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %v", path, err)
		}

		_, err = t.New(path).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", path, err)
		}

		tm.order = append(tm.order, path) // ← Track order
	}

	tm.templates = t
	return nil
}

// RenderTemplate applies the data struct to all loaded templates and combines them
func (tm *TemplateManager) RenderTemplate(data interface{}) (string, error) {
	if tm.templates == nil {
		return "", errors.New("templates not loaded")
	}

	var out strings.Builder

	for _, name := range tm.order { // ← Use the manually preserved order
		tmpl := tm.templates.Lookup(name)
		if tmpl == nil {
			return "", fmt.Errorf("template %s not found", name)
		}

		if err := tmpl.Execute(&out, data); err != nil {
			return "", fmt.Errorf("failed to execute template %s: %v", name, err)
		}
	}

	return out.String(), nil
}

// Now sends an email using SMTP. Pass nil attachments to send only the HTML body.
func Now(config SMTPConfig, r Recipients, subject, msg string, attachments []Attachment, from ...string) error {
	if config.Username == "" || config.Password == "" || config.Server == "" {
		return errors.New("SMTP configuration is missing required fields")
	}

	if config.Port == 0 {
		config.Port = 587
	}

	if len(r.To) == 0 {
		return errors.New("no primary recipient specified")
	}

	finalFrom := config.From
	if len(from) > 0 && from[0] != "" {
		finalFrom = from[0]
	} else if finalFrom == "" {
		finalFrom = config.Username
	}

	m := gomail.NewMessage()
	m.SetHeader("From", finalFrom)
	m.SetHeader("To", r.To...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	if len(r.Cc) > 0 {
		m.SetHeader("Cc", r.Cc...)
	}
	if len(r.Bcc) > 0 {
		m.SetHeader("Bcc", r.Bcc...)
	}

	if err := addAttachments(m, attachments); err != nil {
		return err
	}

	d := gomail.NewDialer(config.Server, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{
		ServerName:         config.Server,
		InsecureSkipVerify: config.Mode == Test || config.Mode == Debug,
	}

	if config.Mode == Live || config.Mode == Test {
		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	} else {
		log.Printf("'🛠 [DEBUG MODE] Email not sent. Details:\nFrom: %s\nTo: %+v\nCC: %+v\nBCC: %+v\nSubject: %s\n", finalFrom, r.To, r.Cc, r.Bcc, subject)
	}

	return nil
}

func addAttachments(m *gomail.Message, attachments []Attachment) error {
	contentIDs := make(map[string]bool)
	for _, attachment := range attachments {
		if attachment.Path == "" {
			return errors.New("attachment path is empty")
		}
		if !attachment.Inline && attachment.ContentID != "" {
			return fmt.Errorf("attachment %q: ContentID is only supported for inline attachments", attachment.Path)
		}

		name := filepath.Base(attachment.Path)
		headers := map[string][]string{
			"Content-Disposition": {mime.FormatMediaType("attachment", map[string]string{"filename": name})},
		}
		if attachment.Inline {
			contentID := attachment.ContentID
			if contentID == "" {
				contentID = name
			}
			for _, c := range contentID {
				if c <= ' ' || c >= 127 || c == '<' || c == '>' {
					return fmt.Errorf("attachment %q: invalid inline ContentID %q", attachment.Path, contentID)
				}
			}
			if contentIDs[contentID] {
				return fmt.Errorf("attachment %q: duplicate inline ContentID %q", attachment.Path, contentID)
			}
			contentIDs[contentID] = true
			headers["Content-Disposition"] = []string{mime.FormatMediaType("inline", map[string]string{"filename": name})}
			headers["Content-ID"] = []string{"<" + contentID + ">"}
		}
		contentType := mime.TypeByExtension(filepath.Ext(name))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return fmt.Errorf("attachment %q: invalid content type: %w", attachment.Path, err)
		}
		params["name"] = name
		headers["Content-Type"] = []string{mime.FormatMediaType(mediaType, params)}

		content, err := os.ReadFile(attachment.Path)
		if err != nil {
			return fmt.Errorf("failed to read attachment %q: %w", attachment.Path, err)
		}
		settings := []gomail.FileSetting{
			gomail.SetHeader(headers),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := io.Copy(w, bytes.NewReader(content))
				return err
			}),
		}
		if attachment.Inline {
			m.Embed(attachment.Path, settings...)
		} else {
			m.Attach(attachment.Path, settings...)
		}
	}
	return nil
}
