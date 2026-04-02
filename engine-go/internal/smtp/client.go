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

func newSMTPClientNoAuth(
	ctx context.Context,
	addr string,
	app *app.ApplicationContext,
) (*smtp.Client, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("TCP connect failed: %w", err)
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
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
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
	tlsCfg := &tls.Config{ServerName: app.Config.SMTP.Host, MinVersion: tls.VersionTLS13}

	var netSMTPclient *smtp.Client
	var err error
	switch app.Config.SMTP.TLSType {
	case "TLS":
		netSMTPclient, err = newSMTPClientWithTLS(ctx, addr, tlsCfg, app)
	case "STARTTLS":
		netSMTPclient, err = newSMTPClientWithSTARTTLS(ctx, addr, tlsCfg, app)
	default:
		netSMTPclient, err = newSMTPClientNoAuth(ctx, addr, app)
	}

	if err != nil {
		return nil, err
	}
	if app.Config.SMTP.TLSType != "NONE" {
		auth := smtp.PlainAuth("", app.Config.SMTP.Username, app.Config.SMTP.Password.SafeString(), app.Config.SMTP.Host)
		if err = netSMTPclient.Auth(auth); err != nil {
			_ = netSMTPclient.Close()
			return nil, fmt.Errorf("smtp authentication failed: %w", err)
		}
	}

	return netSMTPclient, nil
}

func TestSMTPConnection(ctx context.Context, app *app.ApplicationContext) error {
	client, err := newSMTPClient(ctx, app)
	if err != nil {
		return err
	}
	return client.Quit()
}
