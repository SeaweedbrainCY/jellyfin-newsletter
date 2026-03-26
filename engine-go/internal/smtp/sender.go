package smtp

import (
	"context"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"go.uber.org/zap"
)

type EmailMIMEData struct {
	From    string
	To      string
	Subject string
	HTML    string
}

func buildMIMEMessage(email EmailMIMEData) []byte {
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

func sendEmail(
	smtpClient *smtp.Client,
	fromEmailAddr, recipientEmailAddr string,
	emailData EmailMIMEData,
) error {
	if err := smtpClient.Mail(fromEmailAddr); err != nil {
		return fmt.Errorf("MAIL FROM: %w. Given value:%s", err, fromEmailAddr)
	}
	if err := smtpClient.Rcpt(recipientEmailAddr); err != nil {
		return fmt.Errorf("RCPT TO: %w. Given value:%s", err, recipientEmailAddr)
	}

	wc, err := smtpClient.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	defer wc.Close()

	_, err = wc.Write(buildMIMEMessage(emailData))
	return err
}

func SendEmailToAllRecipients(emailHTML string, app *app.ApplicationContext) error {
	const emailSendingDelaySeconds = 2
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
		return fmt.Errorf(
			"fatal error while parsing the FROM sender address. Sender address: %s. Error: %w",
			app.Config.SMTP.SenderName,
			err,
		)
	}

	emailData := EmailMIMEData{
		From:    app.Config.SMTP.SenderName,
		Subject: emailSubject,
		HTML:    emailHTML,
	}

	smtpClient, err := newSMTPClient(context.Background(), app)
	if err != nil {
		return err
	}

	for _, recipient := range app.Config.EmailRecipients {
		cleanedRecipientEmailAddr, err := getEmailAddressFromFriendlyName(recipient)
		if err != nil {
			app.Logger.Error(
				"fatal error while parsing the RCPT recipient address.",
				zap.String("Recipient", recipient),
				zap.Error(err),
			)
			continue
		}
		emailData.To = recipient
		err = sendEmail(
			smtpClient,
			cleanedFromEmailAddr,
			cleanedRecipientEmailAddr,
			emailData,
		)
		if err != nil {
			app.Logger.Error("Failed to send email to "+recipient, zap.String("recipient", recipient), zap.Error(err))
		} else {
			app.Logger.Info("Successfully sent newsletter to " + recipient)
		}
		_ = smtpClient.Reset()
		// We avoid SMTP rate limiting
		time.Sleep(emailSendingDelaySeconds * time.Second)
	}
	return nil
}
