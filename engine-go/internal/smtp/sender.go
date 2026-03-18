package smtp

import (
	"fmt"
	"strings"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
)

type EmailData struct {
	From    string
	To      string
	Subject string
	HTML    string
}

func buildMIMEMessage(email EmailData) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&sb, "From: %s\r\n", email.From)
	fmt.Fprintf(&sb, "To: %s\r\n", email.To)
	fmt.Fprintf(&sb, "Subject: %s\r\n", email.Subject)
	fmt.Fprintf(&sb, "Content-Type: text/html; charset=\"UTF-8\"\r\n")
	fmt.Fprintf(&sb, "\r\n")
	fmt.Fprintf(&sb, "%s", email.HTML)
	return []byte(sb.String())
}

func sendEmail(recipient, emailHTML string, app *app.ApplicationContext) error {
	emailSubject, err := template.BuildEmailTitleWithPlaceholders(
		app.Config.EmailTemplate.Subject,
		app.Config.Jellyfin.ObservedPeriodDays,
		app,
	)
	if err != nil {
		return fmt.Errorf("error while building email's subject: %w", err)
	}
	emailData := EmailData{
		From:    app.Config.SMTP.SenderName,
		To:      recipient,
		Subject: emailSubject,
		HTML:    emailHTML,
	}
	client, err := newSMTPClient(app)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Quit()
	}()

	if err = client.Mail(emailData.From); err != nil {
		return fmt.Errorf("MAIL FROM: %w", err)
	}
	if err = client.Rcpt(emailData.To); err != nil {
		return fmt.Errorf("RCPT TO: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write(buildMIMEMessage(emailData))
	return err
}
