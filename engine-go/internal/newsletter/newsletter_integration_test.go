//go:build integration

package newsletter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/i18n"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/newsletter"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/template"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/tmdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

type FakeClock struct{}

type MailpitContainer struct {
	testcontainers.Container
	Host     string
	SMTPPort int
	APIPort  int
}

type MailpitMessage struct {
	ID      string           `json:"ID"`
	Subject string           `json:"Subject"`
	From    MailpitAddress   `json:"From"`
	To      []MailpitAddress `json:"To"`
	Snippet string           `json:"Snippet"` // plain-text preview
}

type MailpitAddress struct {
	Address string `json:"Address"`
	Name    string `json:"Name"`
}

type MailpitMessagesResponse struct {
	Messages []MailpitMessage `json:"messages"`
	Total    int              `json:"total"`
}

func (c FakeClock) Now() time.Time {
	// Fixtures data are fixed in the time. They are currently captured on 2026-04-01.
	// Altering this date or the fixtures without aletring the other one will break the tests.
	return time.Date(2026, 04, 01, 12, 00, 00, 00, time.UTC)
}

func initApp(t *testing.T) (*app.ApplicationContext, *observer.ObservedLogs, error) {
	// We first load a real Config from a test-dedicated YAML file
	// Then we add the data we need to capture elements (logs, http requests ...)
	configPath := os.Getenv("INTEGRATION_TEST_CONFIG_FILE")
	if configPath == "" {
		t.Fatal("INTEGRATION_TEST_CONFIG_FILE is not defined")
	}
	config, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, nil, err
	}
	localizer, _ := i18n.NewLocalizer("en")
	loggerCore, recordedLogs := observer.New(zap.InfoLevel)
	testCore := zaptest.NewLogger(t).Core()
	logger := zap.New(zapcore.NewTee(loggerCore, testCore), zap.WithFatalHook(zapcore.WriteThenGoexit))
	clock := FakeClock{}
	app := app.ApplicationContext{
		Logger:    logger,
		Clock:     clock,
		Config:    config,
		Localizer: localizer,
	}
	return &app, recordedLogs, nil
}

func removeSecretFromAuthorizationHeader(i *cassette.Interaction) error {
	i.Request.Headers.Del("Authorization")
	return nil
}

func customMatcher(r *http.Request, i cassette.Request) bool {
	// Since we scrub some header when writting the fixture, we adapt the match to omitt headers check
	return r.Method == i.Method && r.URL.String() == i.URL
}

func TestJellyfinNewsletter(t *testing.T) {
	jellyfinRec, err := recorder.New(
		"../../testdata/fixtures/jellyfin", recorder.WithMode(recorder.ModeRecordOnce),
		recorder.WithHook(removeSecretFromAuthorizationHeader, recorder.BeforeSaveHook),
		recorder.WithMatcher(customMatcher),
	)

	require.NoError(t, err)

	defer jellyfinRec.Stop()
	jellyfinHTTPClient := &http.Client{
		Transport: jellyfinRec,
	}

	tmdbRec, err := recorder.New(
		"../../testdata/fixtures/tmdb", recorder.WithMode(recorder.ModeRecordOnce),
		recorder.WithHook(removeSecretFromAuthorizationHeader, recorder.BeforeSaveHook),
		recorder.WithMatcher(customMatcher),
	)
	require.NoError(t, err)

	defer tmdbRec.Stop()
	tmdbHTTPClient := &http.Client{
		Transport: tmdbRec,
	}

	app, _, err := initApp(t)
	require.NoError(t, err)

	err = template.CheckIfThemeIsAvailable(app)
	require.NoError(t, err)
	newsletterWorkflow := newsletter.Workflow{
		JellyfinClient: jellyfin.NewJellyfinAPIClient(jellyfinHTTPClient, app),
		TMDBClient:     tmdb.InitTMDBApiClient(tmdbHTTPClient, app),
	}

	mailpitCT, err := StartMailpit(context.Background(), t)
	require.NoError(t, err)
	app.Config.SMTP.Host = mailpitCT.Host
	app.Config.SMTP.Port = mailpitCT.SMTPPort

	newsletterWorkflow.Run(app, nil)

	messages, err := mailpitCT.GetMessages()
	require.NoError(t, err)
	require.Len(t, messages, 1)
	msg := messages[0]
	assert.Equal(t, "[Jellyfin] New movies ans TV shows of April", msg.Subject)
	html, err := mailpitCT.GetMessageHTML(msg.ID)
	require.NoError(t, err)
	assert.Contains(t, html, "Love Actually")
	assert.Contains(t, html, "Added on 2026-03-31")
	assert.Contains(t, html, "No Time to Die")
	assert.Contains(t, html, "The Last of Us")
	assert.Contains(t, html, "Severance: Season 1, Episodes 2-3")
	assert.Contains(
		t,
		html,
		"You are recieving this email because you are using TestJellyfinUser&#39;s Jellyfin server. If you want to stop receiving these emails, you can unsubscribe by notifying unsubscribe@example.com.",
	)
	assert.Contains(
		t,
		html,
		"Developed with ❤️ by <a href=\"https://github.com/SeaweedbrainCY/\" class=\"footer-link\">SeaweedbrainCY</a> and <a href=\"https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors\" class=\"footer-link\">the contributors</a>.",
	)
	assert.Contains(t, html, "Copyright © 2025 Nathan Stchepinsky, licensed under AGPLv3.")
	assert.NotContains(t, html, "House of the Dragon")
}

func StartMailpit(ctx context.Context, t *testing.T) (*MailpitContainer, error) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/axllent/mailpit:latest",
		ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		WaitingFor: wait.ForHTTP("/api/v1/info").
			WithPort("8025/tcp").
			WithStatusCodeMatcher(func(status int) bool {
				return status == http.StatusOK
			}).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start mailpit: %w", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate mailpit: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("mailpit host: %w", err)
	}

	smtpMappedPort, err := container.MappedPort(ctx, "1025/tcp")
	if err != nil {
		return nil, fmt.Errorf("mailpit port: %w", err)
	}

	smtpPort, err := strconv.Atoi(smtpMappedPort.Port())
	if err != nil {
		return nil, fmt.Errorf("mailpit port: %w", err)
	}

	apiMappedPort, err := container.MappedPort(ctx, "8025/tcp")
	if err != nil {
		return nil, fmt.Errorf("mailpit api port: %w", err)
	}

	apiPort, err := strconv.Atoi(apiMappedPort.Port())
	if err != nil {
		return nil, fmt.Errorf("mailpit api port: %w", err)
	}

	return &MailpitContainer{
		Container: container,
		Host:      host,
		SMTPPort:  smtpPort,
		APIPort:   apiPort,
	}, nil
}

func (m *MailpitContainer) GetMessages() ([]MailpitMessage, error) {
	resp, err := http.Get(m.APIURL() + "/api/v1/messages")
	if err != nil {
		return nil, fmt.Errorf("mailpit get messages: %w", err)
	}
	defer resp.Body.Close()

	var result MailpitMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("mailpit decode messages: %w", err)
	}
	return result.Messages, nil
}

// GetMessageHTML fetches the HTML body of a specific message.
func (m *MailpitContainer) GetMessageHTML(id string) (string, error) {
	resp, err := http.Get(m.APIURL() + "/api/v1/message/" + id)
	if err != nil {
		return "", fmt.Errorf("mailpit get html: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		HTML string `json:"HTML"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("mailpit decode html: %w", err)
	}
	return body.HTML, nil
}

func (m *MailpitContainer) APIURL() string {
	return fmt.Sprintf("http://%s:%d", m.Host, m.APIPort)
}
