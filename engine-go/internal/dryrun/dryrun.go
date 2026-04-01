package dryrun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/clock"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/smtp"
	"go.uber.org/zap"
)

type metadataJSON struct {
	Datetime          string
	SMTPTestResult    string
	NewDetectedMovies []jellyfin.MovieItem
	NewDetectedSeries []jellyfin.NewlyAddedSeriesItem
}

func fillFilenameTemplate(filename string, app *app.ApplicationContext) string {
	templateData := struct {
		Datetime string
	}{
		Datetime: app.Clock.Now().Format("2006-01-02T15:04:05Z07:00"),
	}
	tmpl, err := template.New("filename").Option("missingkey=zero").Parse(filename)
	if err != nil {
		app.Logger.Debug(
			"An error occured while filling the dry-run output filename template",
			zap.String("step", "create"),
			zap.String("filename", filename),
			zap.String("datetime", templateData.Datetime),
			zap.Error(err),
		)
		return filename
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		app.Logger.Debug(
			"An error occured while filling the dry-run output filename template",
			zap.String("step", "execute"),
			zap.String("filename", filename),
			zap.String("datetime", templateData.Datetime),
			zap.Error(err),
		)
		return filename
	}

	return buf.String()
}

func marshalNewItems(items any) string {
	marshalledBytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(marshalledBytes)
}

func addMetadataToHTML(emailHTML string, newJellyfinMovies *[]jellyfin.MovieItem,
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem, smtpTestResult string, clock clock.ClockInterface) string {
	if smtpTestResult == "" {
		smtpTestResult = "Not tested"
	}

	metadata := fmt.Sprintf(
		"<!--\nJellyfin-newsletter dry run\nGenerated at: %s\nSMTP test result:%s\nNew movies detected: %s\nNew series detected: %s\n-->\n\n",
		clock.Now().Format("2006-01-02T15:04:05Z07:00"),
		smtpTestResult,
		marshalNewItems(newJellyfinMovies),
		marshalNewItems(newJellyfinSeries),
	)
	return metadata + emailHTML
}

func saveMetadataAsJSONFile(outputDirectory, outputFilename string, newJellyfinMovies *[]jellyfin.MovieItem,
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem, smtpTestResult string, clock clock.ClockInterface) error {
	metadata := metadataJSON{
		NewDetectedMovies: *newJellyfinMovies,
		NewDetectedSeries: *newJellyfinSeries,
		Datetime:          clock.Now().Format("2006-01-02T15:04:05Z07:00"),
		SMTPTestResult:    smtpTestResult,
	}

	filePath := filepath.Join(outputDirectory, outputFilename)
	//nolint:musttag // It's debugging data, data consistency is not important
	metadataMarshalled, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, metadataMarshalled, 0600)
	return err
}

func saveHTMLFile(outputDirectory, outputFilename, emailHTML string) error {
	filePath := filepath.Join(outputDirectory, outputFilename)
	return os.WriteFile(filePath, []byte(emailHTML), 0600)
}

func SaveDryRunEmail(emailHTML string, newJellyfinMovies *[]jellyfin.MovieItem,
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem, app *app.ApplicationContext) {
	outputFilename := fillFilenameTemplate(app.Config.DryRun.OutputFilename, app)
	smtpTestResult := "SMTP connection not tested."
	if app.Config.DryRun.TestSMTPConnection {
		app.Logger.Info("Testing SMTP connection ...")
		smtpTestErr := smtp.TestSMTPConnection(context.Background(), app)
		if smtpTestErr != nil {
			app.Logger.Warn("SMTP connection FAILED. No email sent.", zap.Error(smtpTestErr))
			smtpTestResult = "SMTP connection FAILED. No email sent. Error: " + smtpTestErr.Error()
		} else {
			app.Logger.Info("SMTP connection SUCCESSFUL. No email sent.")
			smtpTestResult = "SMTP connection SUCCESSFUL. No email sent."
		}
	}

	if app.Config.DryRun.IncludeMetadata {
		emailHTML = addMetadataToHTML(emailHTML, newJellyfinMovies, newJellyfinSeries, smtpTestResult, app.Clock)
	}

	if app.Config.DryRun.SaveEmailData {
		filename := strings.ReplaceAll(outputFilename, ".html", ".json")
		err := saveMetadataAsJSONFile(
			app.Config.DryRun.OutputDirectory,
			filename,
			newJellyfinMovies,
			newJellyfinSeries,
			smtpTestResult,
			app.Clock,
		)
		if err != nil {
			app.Logger.Error(
				"Impossible to write metadata in json file.",
				zap.String("output directory", app.Config.DryRun.OutputDirectory),
				zap.String("filename", filename),
				zap.Error(err),
			)
		}
	}

	err := saveHTMLFile(app.Config.DryRun.OutputDirectory, outputFilename, emailHTML)
	if err != nil {
		app.Logger.Error(
			"An error occurred while saving the HTML email file.",
			zap.String("output directory", app.Config.DryRun.OutputDirectory),
			zap.String("filename", outputFilename),
			zap.Error(err),
		)
	}
}
