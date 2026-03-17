package dryrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/jellyfin"
)

func fillFilenameTemplate(filename string) string {
	templateData := struct {
		datetime string
	}{
		datetime: time.Now().Format("RFC3339"),
	}
	tmpl, err := template.New("filename").Option("missingkey=zero").Parse(filename)
	if err != nil {
		return filename
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
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
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem, SMTPTestResult string) string {
	if SMTPTestResult == "" {
		SMTPTestResult = "Not tested"
	}

	metadata := fmt.Sprintf("<!--\nJellyfin-newsletter dry run\nGenerated at: %s\nSMTP test result:%s\nNew movies detected: %s\nNew series detected: %s\n-->\n\n", time.Now().Format("RFC3339"), SMTPTestResult, marshalNewItems(newJellyfinMovies), marshalNewItems(newJellyfinSeries))
	return metadata + emailHTML
}

func SaveDryRunEmail(emailHTML string, newJellyfinMovies *[]jellyfin.MovieItem,
	newJellyfinSeries *[]jellyfin.NewlyAddedSeriesItem, app *app.ApplicationContext) {
	outputDirectory := "/app/config/previews/"
	outputFilenameTemplate := "newsletter_{{.datetime}}.html"

	if app.Config.DryRun.OutputDirectory == "" {
		outputDirectory = app.Config.DryRun.OutputDirectory
	}

	if app.Config.DryRun.OutputFilename == "" {
		outputFilenameTemplate = app.Config.DryRun.OutputFilename
	}

	outputFilename := fillFilenameTemplate(outputFilenameTemplate)

	if app.Config.DryRun.IncludeMetadata {
		emailHTML = addMetadataToHTML(emailHTML, newJellyfinMovies, newJellyfinSeries)
	}
}
