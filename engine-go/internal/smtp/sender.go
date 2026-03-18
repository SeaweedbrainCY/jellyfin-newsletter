package smtp

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"go.uber.org/zap"
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

// SMTP MAIL FROM command expects the email to be a@b.c and doesn't access a <a@b.c> or <a@b.c>. This utility extract a@b.c from a <a@b.c> or <a@b.c>. If not found, it returns an error.
func getEmailAddressFromFriendlyName(emailFriendlyName string) (string, error) {
	parsedAddress, err := mail.ParseAddress(emailFriendlyName)
	if err != nil {
		return "", err
	}
	return parsedAddress.Address, nil
}

func sendEmail(ctx context.Context, recipient, emailHTML string, app *app.ApplicationContext) error {
	emailSubject, err := template.BuildEmailTitleWithPlaceholders(
		app.Config.EmailTemplate.Subject,
		app.Config.Jellyfin.ObservedPeriodDays,
		app,
	)
	if err != nil {
		return fmt.Errorf("error while building email's subject: %w", err)
	}

	cleanedFromEmailAddr, err := getEmailAddressFromFriendlyName(app.Config.SMTP.SenderName)
	if err != nil {
		return fmt.Errorf("Fatal error while parsing the FROM sender address. Sender address: %s. Error: %w.", app.Config.SMTP.SenderName, err)
	}

	cleanedRecipientEmailAddr, err := getEmailAddressFromFriendlyName(recipient)
	if err != nil {
		return fmt.Errorf("Fatal error while parsing the RCPT recipient address. Recipient address: %s. Error: %w.", recipient, err)
	}

	client, err := newSMTPClient(ctx, app)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Quit()
	}()

	if err = client.Mail(cleanedFromEmailAddr); err != nil {
		return fmt.Errorf("MAIL FROM: %w. Given value:%s", err, cleanedFromEmailAddr)
	}
	if err = client.Rcpt(cleanedRecipientEmailAddr); err != nil {
		return fmt.Errorf("RCPT TO: %w. Given value:%s", err, cleanedRecipientEmailAddr)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	defer wc.Close()

	emailData := EmailData{
		From:    app.Config.SMTP.SenderName,
		To:      recipient,
		Subject: emailSubject,
		HTML:    emailHTML,
	}

	_, err = wc.Write(buildMIMEMessage(emailData))
	return err
}

func SendEmailToAllRecipients(emailHTML string, app *app.ApplicationContext) {
	for _, recipient := range app.Config.EmailRecipients {
		err := sendEmail(context.Background(), recipient, emailHTML, app)
		if err != nil {
			app.Logger.Error("Failed to send email to "+recipient, zap.String("recipient", recipient), zap.Error(err))
		} else {
			app.Logger.Info("Successfully sent newsletter to " + recipient)
		}
	}
}
