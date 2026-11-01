# GoSend Email Package (v3)

<!-- Badges -->
![GitHub Repo stars](https://img.shields.io/github/stars/aro-wolo/gosend?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/aro-wolo/gosend)
![GitHub issues](https://img.shields.io/github/issues/aro-wolo/gosend)
![GitGitHub contributors](https://img.shields.io/github/contributors/aro-wolo/gosend)
[![Go Report Card](https://goreportcard.com/badge/github.com/aro-wolo/gosend)](https://goreportcard.com/report/github.com/aro-wolo/gosend)
[![Build Status](https://github.com/aro-wolo/gosend/actions/workflows/go.yml/badge.svg)](https://github.com/aro-wolo/gosend/actions)
[![codecov](https://codecov.io/gh/aro-wolo/gosend/branch/main/graph/badge.svg)](https://codecov.io/gh/aro-wolo/gosend)
![Go Module](https://img.shields.io/github/go-mod/go-version/aro-wolo/gosend)
[![License](https://img.shields.io/github/license/aro-wolo/gosend.svg)](https://github.com/aro-wolo/gosend/blob/main/LICENSE)

---

## Overview

**GoSend** is a flexible and developer‑friendly email‑sending package for Go.

It supports:

- SMTP email sending  
- Multiple environments: Debug, Test, Live  
- HTML email bodies  
- Normal and inline file attachments
- Template-based emails  
- Multi-file template stacking (header/body/footer)  
- TLS configuration  
- High-level helper API (`SendMail`) — *new in v2*  
- Low-level API (`Now`)  

GoSend is lightweight, dependency‑minimal, and easy to integrate into any Go project.

---

# What's New in v3.0.0

- Send normal attachments or inline files with HTML `cid:` references.
- Pass attachments directly to `Now` and `SendMail`, separately from SMTP settings.
- Validate attachment paths and inline Content-IDs before delivery.

## Migrating from v2

This release changes both sending function signatures. Update imports to
`github.com/aro-wolo/gosend/v3` and add the attachments argument (or `nil`):

```go
// No attachments, with an optional sender override.
err := gosend.Now(config, recipients, "Hello", "<p>Hello!</p>",
    nil, "sender@example.com")

// No attachments, with an optional template directory.
err = gosend.SendMail(config, recipients, "Welcome", "welcome",
    data, nil, "templates")
```

The attachment argument follows the HTML body in `Now`, and the template data
in `SendMail`. Optional sender and template directory arguments stay last.
`SMTPConfig` remains unchanged.

## Features introduced in v2

### ✔ Added High-Level Helper: `SendMail()`
A simple wrapper to send template emails in one line.

### ✔ Optional Template Base Path
Allows custom template folder structures.

### ✔ Standardized Template Layout
Header → Body → Footer.

---

## Installation

```bash
go get github.com/aro-wolo/gosend/v3@v3.0.0
```

---

# 📦 Usage

## 1. Basic Email (No Templates)

```go
config := gosend.SMTPConfig{
    Username: "your-email@example.com",
    Password: "your-password",
    Server:   "smtp.example.com",
    Port:     587,
    Mode:     gosend.Live,
}

recipients := gosend.Recipients{
    To: []string{"user@example.com"},
}

err := gosend.Now(config, recipients, "Hello", "<p>Hi there!</p>", nil)
if err != nil {
    log.Fatal(err)
}
```

---

# 📨 Template Email (Recommended)

### Folder Structure

```
templates/
  header.html
  welcome.html
  footer.html
```

### Sample Template Usage

```go
err := gosend.SendMail(
    config,
    gosend.Recipients{To: []string{"user@example.com"}},
    "Welcome!",
    "welcome", // loads welcome.html
    map[string]any{
        "Name": "John Doe",
    },
    nil, // no attachments
    "templates", // optional base path
)
```

---

# 🧩 High-Level Helper API

```go
func SendMail(
    config SMTPConfig,
    recipients Recipients,
    subject string,
    templateName string,
    data map[string]any,
    attachments []Attachment,
    templateBasePath ...string,
) error
```

### Default Template Resolution

```
base/header.html
base/{templateName}.html
base/footer.html
```

### Example

```
templates/header.html
templates/verify.html
templates/footer.html
```

Usage:

```go
gosend.SendMail(config, r, "Verify Email", "verify", data, nil, "templates")
```

---

## Attachments

Pass an `[]gosend.Attachment` parameter to either `Now` (after the HTML body)
or `SendMail` (after the template data). Attachments are per-email arguments,
not SMTP configuration. A file is a normal, downloadable attachment unless
`Inline` is `true`. Both types can be used in
the same email.

```go
attachments := []gosend.Attachment{
    {Path: "documents/report.pdf"},
    {Path: "images/logo.png", Inline: true, ContentID: "company-logo"},
}

err := gosend.Now(config, recipients, "Your report",
    `<p>Your report is attached.</p><img src="cid:company-logo" alt="Logo">`,
    attachments)
if err != nil {
    log.Fatal(err)
}
```

For template emails, put the same `cid:company-logo` reference in the HTML
template and pass the attachments to `SendMail`:

```go
err := gosend.SendMail(config, recipients, "Your report", "report",
    map[string]any{"Name": "John"}, attachments, "templates")
if err != nil {
    log.Fatal(err)
}
```

| Attachment field | Description |
|------------------|-------------|
| Path | Required path to a readable file on disk |
| Inline | `false` (default) for a normal attachment; `true` for an inline MIME part |
| ContentID | Optional inline identifier without `<` or `>`; defaults to the file's base name (e.g. `logo.png`, referenced as `cid:logo.png`) |

Inline Content-IDs must be unique within the email and contain only printable
ASCII characters, excluding spaces and angle brackets. Use an explicit
Content-ID when the file name contains spaces or non-ASCII characters.
Setting `ContentID` on a normal attachment returns an error.

Files are read into memory before delivery, including in Debug mode. Missing
or unreadable files and invalid attachment settings return errors; Debug still
does not send any email. File names and MIME types are derived from the paths.
Pass an empty slice or `nil` for no attachments. Existing calls must add `nil`
after the HTML body (`Now`) or template data (`SendMail`). Optional sender and
template base path arguments remain last.

---

# ⚙ Configuration Reference

### SMTPConfig Fields

| Field | Description |
|-------|-------------|
| Username | SMTP username |
| Password | SMTP password |
| Server | SMTP hostname |
| Port | SMTP port (default: 587) |
| Mode | Debug / Test / Live |
| From | Optional From header |

### Environment Modes

| Mode | Behavior |
|------|----------|
| Debug | Logs email instead of sending |
| Test | Sends email without strict TLS |
| Live | Full SMTP delivery |

---

# 🧪 Low-Level Template API

### Now

```go
func Now(
    config SMTPConfig,
    recipients Recipients,
    subject, msg string,
    attachments []Attachment,
    from ...string,
) error
```

### TemplateManager

```go
tm := gosend.NewTemplateManager()
tm.ParseTemplate("header.html", "body.html", "footer.html")
result, err := tm.RenderTemplate(data)
```


# 🤝 Contributing

Pull requests are welcome!  
Open issues here:  
https://github.com/aro-wolo/gosend/issues

---

# 📄 License

MIT License.
