package smtp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
)

func newSMTPClientWithTLS(
	ctx context.Context,
	addr string,
	tlsCfg *tls.Config,
	app *app.ApplicationContext,
) (*smtp.Client, error) {
	dialer := &tls.Dialer{NetDialer: &net.Dialer{}, Config: tlsCfg}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("TLS connect failed: %w", err)
	}
	client, err := smtp.NewClient(conn, app.Config.SMTP.Host)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("SMTP handshake failed: %w", err)
	}
	return client, nil
}

func newSMTPClientWithSTARTTLS(
	ctx context.Context,
	addr string,
	tlsCfg *tls.Config,
	app *app.ApplicationContext,
) (*smtp.Client, error) {
	dialer := &tls.Dialer{NetDialer: &net.Dialer{}, Config: tlsCfg}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("TLS connect failed: %w", err)
	}
	client, err := smtp.NewClient(conn, app.Config.SMTP.Host)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("SMTP handshake failed: %w", err)
	}
	ok, _ := client.Extension("STARTTLS")
	if !ok {
		_ = client.Close()
		return nil, errors.New("server does not advertise STARTTLS")
	}
	if err = client.StartTLS(tlsCfg); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("STARTTLS handshake failed: %w", err)
	}
	return client, nil
}

func newSMTPClient(ctx context.Context, app *app.ApplicationContext) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", app.Config.SMTP.Host, app.Config.SMTP.Port)
	auth := smtp.PlainAuth("", app.Config.SMTP.Username, app.Config.SMTP.Password.SafeString(), app.Config.SMTP.Host)
	tlsCfg := &tls.Config{ServerName: app.Config.SMTP.Host, MinVersion: tls.VersionTLS13}

	var client *smtp.Client
	var err error
	if app.Config.SMTP.TLSType == "TLS" {
		client, err = newSMTPClientWithTLS(ctx, addr, tlsCfg, app)
	} else {
		client, err = newSMTPClientWithSTARTTLS(ctx, addr, tlsCfg, app)
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

func TestSMTPConnection(ctx context.Context, app *app.ApplicationContext) error {
	client, err := newSMTPClient(ctx, app)
	if err != nil {
		return err
	}
	return client.Quit()
}
