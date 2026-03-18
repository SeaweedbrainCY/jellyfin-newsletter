package smtp

import (
	"fmt"
	"strings"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
)

type emailData struct {
	From    string
	To      string
	Subject string
	HTML    string
}

func buildMIMEMessage(email emailData) []byte {
	var sb strings.Builder

	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	sb.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	sb.WriteString("\r\n") // blank line separates headers from body
	sb.WriteString(email.HTML)

	return []byte(sb.String())
}

func sendEmail(recipient, emailHTML string, app *app.ApplicationContext) error {
	emailSubject, err := template.BuildEmailTitleWithPlaceholders(app.Config.EmailTemplate.Subject, app.Config.Jellyfin.ObservedPeriodDays, app)
	if err != nil {
		return fmt.Errorf("Error while building email's subject: %w", err)
	}
	emailData := emailData{
		From:    app.Config.SMTP.SenderName,
		To:      recipient,
		Subject: emailSubject,
		HTML:    emailHTML,
	}
	client, err := newSMTPClient(app)
	if err != nil {
		return err
	}
	defer client.Quit()

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
