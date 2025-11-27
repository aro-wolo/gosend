# GoSend Email Package (v2)

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
- Template-based emails  
- Multi-file template stacking (header/body/footer)  
- TLS configuration  
- High-level helper API (`SendMail`) — *new in v2*  
- Low-level API (`Now`)  

GoSend is lightweight, dependency‑minimal, and easy to integrate into any Go project.

---

# 🚀 What’s New in v2

### ✔ Added High-Level Helper: `SendMail()`
A simple wrapper to send template emails in one line.

### ✔ Optional Template Base Path
Allows custom template folder structures.

### ✔ Standardized Template Layout
Header → Body → Footer.

---

## Installation

```bash
go get github.com/aro-wolo/gosend/v2
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

err := gosend.Now(config, recipients, "Hello", "<p>Hi there!</p>")
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
    "templates", // optional base path
)
```

---

# 🧩 High-Level Helper API (v2)

```go
func SendMail(
    config SMTPConfig,
    recipients Recipients,
    subject string,
    templateName string,
    data map[string]any,
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
gosend.SendMail(config, r, "Verify Email", "verify", data, "templates")
```

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

