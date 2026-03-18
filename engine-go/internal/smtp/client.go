package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
)

func newSMTPClientWithTLS(addr string, auth smtp.Auth, tlsCfg *tls.Config, app *app.ApplicationContext) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return nil, fmt.Errorf("TLS connect failed: %w", err)
	}
	client, err := smtp.NewClient(conn, app.Config.SMTP.Host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("SMTP handshake failed: %w", err)
	}
	return client, nil
}

func newSMTPClientWithSTARTTLS(addr string, auth smtp.Auth, tlsCfg *tls.Config, app *app.ApplicationContext) (*smtp.Client, error) {
	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("TCP connect failed: %w", err)
	}
	ok, _ := client.Extension("STARTTLS")
	if !ok {
		_ = client.Close()
		return nil, fmt.Errorf("server does not advertise STARTTLS")
	}
	if err = client.StartTLS(tlsCfg); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("STARTTLS handshake failed: %w", err)
	}
	return client, nil
}

func newSMTPClient(app *app.ApplicationContext) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", app.Config.SMTP.Host, app.Config.SMTP.Port)
	auth := smtp.PlainAuth("", app.Config.SMTP.Username, app.Config.SMTP.Password.SafeString(), app.Config.SMTP.Host)
	tlsCfg := &tls.Config{ServerName: app.Config.SMTP.Host}

	var client *smtp.Client
	var err error
	if app.Config.SMTP.TLSType == "TLS" {
		client, err = newSMTPClientWithTLS(addr, auth, tlsCfg, app)
	} else {
		client, err = newSMTPClientWithSTARTTLS(addr, auth, tlsCfg, app)
	}

	if err != nil {
		return nil, err
	}

	if err = client.Auth(auth); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return client, nil
}

func TestSMTPConnection(app *app.ApplicationContext) error {
	client, err := newSMTPClient(app)
	if err != nil {
		return err
	}
	return client.Quit()
}
