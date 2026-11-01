package gosend

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/gomail.v2"
)

type mimePart struct {
	header textproto.MIMEHeader
	body   []byte
}

func readMIMEParts(t *testing.T, header textproto.MIMEHeader, body io.Reader) []mimePart {
	t.Helper()
	contentType, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(contentType, "multipart/") {
		var parts []mimePart
		reader := multipart.NewReader(body, params["boundary"])
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				return parts
			}
			if err != nil {
				t.Fatal(err)
			}
			parts = append(parts, readMIMEParts(t, part.Header, part)...)
		}
	}
	if header.Get("Content-Transfer-Encoding") == "base64" {
		body = base64.NewDecoder(base64.StdEncoding, body)
	} else if header.Get("Content-Transfer-Encoding") == "quoted-printable" {
		body = quotedprintable.NewReader(body)
	}
	content, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}
	return []mimePart{{header: header, body: content}}
}

func attachmentFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAttachmentMIME(t *testing.T) {
	dir := t.TempDir()
	payload := []byte{0, 1, 2, 127, 128, 255}
	path := attachmentFile(t, dir, "logo.png", payload)
	quotedPath := attachmentFile(t, dir, `report "final".unknown`, payload)
	for _, tc := range []struct {
		name        string
		attachments []Attachment
		rootType    string
		disposition []string
		contentIDs  []string
	}{
		{"none", nil, "text/html", nil, nil},
		{"normal", []Attachment{{Path: path}}, "multipart/mixed", []string{"attachment"}, []string{""}},
		{"inline default", []Attachment{{Path: path, Inline: true}}, "multipart/related", []string{"inline"}, []string{"<logo.png>"}},
		{"inline custom", []Attachment{{Path: path, Inline: true, ContentID: "logo@example.com"}}, "multipart/related", []string{"inline"}, []string{"<logo@example.com>"}},
		{"mixed", []Attachment{{Path: quotedPath}, {Path: path, Inline: true, ContentID: "logo"}}, "multipart/mixed", []string{"inline", "attachment"}, []string{"<logo>", ""}},
		{"multiple normal", []Attachment{{Path: path}, {Path: quotedPath}}, "multipart/mixed", []string{"attachment", "attachment"}, []string{"", ""}},
		{"multiple inline", []Attachment{{Path: path, Inline: true, ContentID: "one"}, {Path: path, Inline: true, ContentID: "two"}}, "multipart/related", []string{"inline", "inline"}, []string{"<one>", "<two>"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := gomail.NewMessage()
			m.SetBody("text/html", `<img src="cid:logo">`)
			if err := addAttachments(m, tc.attachments); err != nil {
				t.Fatal(err)
			}
			var raw bytes.Buffer
			if _, err := m.WriteTo(&raw); err != nil {
				t.Fatal(err)
			}
			message, err := mail.ReadMessage(&raw)
			if err != nil {
				t.Fatal(err)
			}
			rootType, _, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
			if err != nil || rootType != tc.rootType {
				t.Fatalf("root MIME type = %q, err = %v; want %q", rootType, err, tc.rootType)
			}
			parts := readMIMEParts(t, textproto.MIMEHeader(message.Header), message.Body)
			if len(parts) != 1+len(tc.attachments) {
				t.Fatalf("got %d MIME parts, want %d", len(parts), 1+len(tc.attachments))
			}
			if string(parts[0].body) != `<img src="cid:logo">` {
				t.Fatalf("HTML body changed: %q", parts[0].body)
			}
			for i, part := range parts[1:] {
				disposition, params, err := mime.ParseMediaType(part.header.Get("Content-Disposition"))
				if err != nil || disposition != tc.disposition[i] {
					t.Fatalf("disposition = %q, err = %v; want %q", disposition, err, tc.disposition[i])
				}
				wantName := "logo.png"
				if tc.name == "mixed" && i == 1 || tc.name == "multiple normal" && i == 1 {
					wantName = filepath.Base(quotedPath)
				}
				if params["filename"] != wantName {
					t.Errorf("filename = %q, want %q", params["filename"], wantName)
				}
				if got := part.header.Get("Content-ID"); got != tc.contentIDs[i] {
					t.Errorf("Content-ID = %q, want %q", got, tc.contentIDs[i])
				}
				mediaType, mediaParams, err := mime.ParseMediaType(part.header.Get("Content-Type"))
				wantType := "image/png"
				if wantName == filepath.Base(quotedPath) {
					wantType = "application/octet-stream"
				}
				if err != nil || mediaType != wantType || mediaParams["name"] != wantName {
					t.Errorf("unexpected content type: %q (%v)", part.header.Get("Content-Type"), err)
				}
				if !bytes.Equal(part.body, payload) {
					t.Errorf("attachment bytes = %v, want %v", part.body, payload)
				}
			}
		})
	}
}

func TestNow_InvalidAttachments(t *testing.T) {
	dir := t.TempDir()
	path := attachmentFile(t, dir, "logo.png", []byte("image"))
	for _, tc := range []struct {
		name        string
		attachments []Attachment
		want        string
	}{
		{"empty path", []Attachment{{}}, "attachment path is empty"},
		{"missing file", []Attachment{{Path: filepath.Join(dir, "missing")}}, "failed to read attachment"},
		{"directory", []Attachment{{Path: dir}}, "failed to read attachment"},
		{"normal CID", []Attachment{{Path: path, ContentID: "logo"}}, "only supported for inline"},
		{"header injection", []Attachment{{Path: path, Inline: true, ContentID: "logo\r\nX-Test: value"}}, "invalid inline ContentID"},
		{"bracketed CID", []Attachment{{Path: path, Inline: true, ContentID: "<logo>"}}, "invalid inline ContentID"},
		{"space in CID", []Attachment{{Path: path, Inline: true, ContentID: "my logo"}}, "invalid inline ContentID"},
		{"non ASCII CID", []Attachment{{Path: path, Inline: true, ContentID: "logo\u00e9"}}, "invalid inline ContentID"},
		{"duplicate CID", []Attachment{{Path: path, Inline: true}, {Path: path, Inline: true}}, "duplicate inline ContentID"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := testSMTPConfig
			err := Now(config, testRecipients, "Subject", "Body", tc.attachments)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
			if tc.name == "missing file" && !errors.Is(err, os.ErrNotExist) {
				t.Errorf("file error is not wrapped: %v", err)
			}
		})
	}
}

func TestAttachmentsDebug(t *testing.T) {
	dir := t.TempDir()
	config := testSMTPConfig
	attachments := []Attachment{
		{Path: attachmentFile(t, dir, "report.txt", []byte("report"))},
		{Path: attachmentFile(t, dir, "logo.png", []byte("image")), Inline: true},
	}
	if err := Now(config, testRecipients, "Attachments", `<img src="cid:logo.png">`, attachments); err != nil {
		t.Fatal(err)
	}

	attachmentFile(t, dir, "header.html", nil)
	attachmentFile(t, dir, "welcome.html", []byte(`<img src="cid:logo.png">`))
	attachmentFile(t, dir, "footer.html", nil)
	attachments[0].Path = filepath.Join(dir, "missing")
	err := SendMail(config, testRecipients, "Attachments", "welcome", nil, attachments, dir)
	if err == nil || !strings.Contains(err.Error(), "failed sending email: failed to read attachment") || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected wrapped attachment error from SendMail, got %v", err)
	}
}

type smtpResult struct {
	raw []byte
	err error
}

func receiveSMTP(conn net.Conn) smtpResult {
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return smtpResult{err: err}
	}
	protocol := textproto.NewConn(conn)
	if err := protocol.PrintfLine("220 localhost ESMTP"); err != nil {
		return smtpResult{err: err}
	}
	var raw []byte
	for {
		line, err := protocol.ReadLine()
		if err != nil {
			return smtpResult{err: err}
		}
		switch {
		case strings.HasPrefix(line, "EHLO "), strings.HasPrefix(line, "HELO "),
			strings.HasPrefix(line, "MAIL FROM:"), strings.HasPrefix(line, "RCPT TO:"):
			err = protocol.PrintfLine("250 OK")
		case line == "DATA":
			if err = protocol.PrintfLine("354 Send message"); err == nil {
				raw, err = io.ReadAll(protocol.DotReader())
			}
			if err == nil {
				err = protocol.PrintfLine("250 Accepted")
			}
		case line == "QUIT":
			err = protocol.PrintfLine("221 Bye")
			return smtpResult{raw: raw, err: err}
		default:
			return smtpResult{err: fmt.Errorf("unexpected SMTP command %q", line)}
		}
		if err != nil {
			return smtpResult{err: err}
		}
	}
}

func TestSendAttachmentsSMTP(t *testing.T) {
	for _, api := range []string{"Now", "SendMail"} {
		t.Run(api, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			result := make(chan smtpResult, 1)
			go func() {
				conn, err := listener.Accept()
				if err != nil {
					result <- smtpResult{err: err}
					return
				}
				result <- receiveSMTP(conn)
			}()
			host, portString, err := net.SplitHostPort(listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			port, err := strconv.Atoi(portString)
			if err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			config := testSMTPConfig
			config.Server, config.Port, config.Mode = host, port, Test
			attachments := []Attachment{
				{Path: attachmentFile(t, dir, "report.txt", []byte("report data"))},
				{Path: attachmentFile(t, dir, "logo.png", []byte("image data")), Inline: true, ContentID: "logo"},
			}
			body := `<p>Hello Tester</p><img src="cid:logo">`
			if api == "Now" {
				err = Now(config, testRecipients, "Attachments", body, attachments, "override@example.com")
			} else {
				attachmentFile(t, dir, "header.html", []byte("<p>Hello {{.Name}}</p>"))
				attachmentFile(t, dir, "welcome.html", []byte(`<img src="cid:logo">`))
				attachmentFile(t, dir, "footer.html", nil)
				err = SendMail(config, testRecipients, "Attachments", "welcome", map[string]any{"Name": "Tester"}, attachments, dir)
			}
			if err != nil {
				t.Fatal(err)
			}
			var received smtpResult
			select {
			case received = <-result:
			case <-time.After(6 * time.Second):
				t.Fatal("timed out waiting for SMTP message")
			}
			if received.err != nil {
				t.Fatal(received.err)
			}
			message, err := mail.ReadMessage(bytes.NewReader(received.raw))
			if err != nil {
				t.Fatal(err)
			}
			wantFrom := config.From
			if api == "Now" {
				wantFrom = "override@example.com"
			}
			if message.Header.Get("From") != wantFrom || message.Header.Get("Subject") != "Attachments" {
				t.Fatalf("unexpected message headers: %v", message.Header)
			}
			parts := readMIMEParts(t, textproto.MIMEHeader(message.Header), message.Body)
			if len(parts) != 3 {
				t.Fatalf("got %d parts, want HTML plus two attachments", len(parts))
			}
			if string(parts[0].body) != body || string(parts[1].body) != "image data" || string(parts[2].body) != "report data" {
				t.Fatalf("unexpected delivered bodies: %v", parts)
			}
			if parts[1].header.Get("Content-ID") != "<logo>" ||
				!strings.HasPrefix(parts[1].header.Get("Content-Disposition"), "inline;") ||
				!strings.HasPrefix(parts[2].header.Get("Content-Disposition"), "attachment;") {
				t.Fatalf("unexpected delivered attachment headers: %v", parts)
			}
		})
	}
}
