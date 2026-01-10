package context

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validConfigYAML = `
scheduler:
  cron: "0 8 1 * *"
    
log:
  level: INFO
  format: console

jellyfin:
  url: http://localhost:8096
  api_token: secret
  watched_film_folders:
    - /movies
  watched_tv_folders:
    - /series
  observed_period_days: 30
  ignore_item_added_before_last_newsletter: false

tmdb:
  api_key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30

email:
  smtp_server: smtp.example.com
  smtp_port: 587
  smtp_username: user
  smtp_password: pass
  smtp_sender_email: Jellyfin
  smtp_tls_type: "TLS"

email_template:
  theme: "classic"
  language: en
  subject: New releases
  title: Newsletter
  subtitle: This week
  jellyfin_url: http://localhost:8096
  unsubscribe_email: unsub@example.com
  jellyfin_owner_name: Admin
  sort_mode: "date_asc"
  display_overview_max_items: 10

dry-run:
  enabled: true
  test_smtp_connection: false
  output_directory: "/app/config/previews/"
  output_filename: "newsletter_{date}.html"
  include_metadata: true
  save_email_data: true

recipients:
  - "user1@example.com"
  - "user2@example.com"
`

// Recursive function
// Will gradually construct the newYaml and parse the base Yaml line by line.
// The pattern will be search by exploring the indentation if the base field is equaled to the fieldPath part.
//   - newYaml represents the final YAML string without the part to be removed
//   - baseYamlLines represents the base YAML file being processed as an array of strings. One string for each line
//   - linePositionToParse represents the line position in baseYamlLines to parse in the current iteration
//   - fieldLocation represents the dot notation of the field/section to remove. For example jellyfin or jellyfin.url
//   - indentationLevelToDelete tells if the current line should be deleted because it has the indentation of a to be deleted section
//   - ignoreIndentationLevel tells if we should be looking at the field in a section with this identation level or wait to go down to search pattern again. For example, if we are exploring email section, ignoreIndentationLevel will be equaled to 1 and search will be paused until ignoreIndentationLevel goes down to 0 again. If equaled to 0, this is ignored
func computeNewYamlAfterPartRemoval(newYaml *string, baseYamlLines *[]string, linePositionToParse int, fieldPath *[]string, indentationLevelToDelete int, ignoreIndentationLevel int) error {
	if linePositionToParse == len(*baseYamlLines) {
		return nil
	} else {
		line := (*baseYamlLines)[linePositionToParse]
		lineTrimmed := strings.TrimLeft(line, " ")
		currentIdentNumber := (len(line) - len(lineTrimmed)) / 2 // 1 indent = 2 spaces
		if currentIdentNumber >= indentationLevelToDelete && indentationLevelToDelete != 0 {
			return computeNewYamlAfterPartRemoval(newYaml, baseYamlLines, linePositionToParse+1, fieldPath, indentationLevelToDelete, 0)
		} else if currentIdentNumber >= ignoreIndentationLevel && ignoreIndentationLevel != 0 {
			*newYaml = *newYaml + "\n" + line
			return computeNewYamlAfterPartRemoval(newYaml, baseYamlLines, linePositionToParse+1, fieldPath, 0, ignoreIndentationLevel)
		} else if currentIdentNumber > len(*fieldPath)-1 {
			// We are too high in indent, can be ignored
			*newYaml = *newYaml + "\n" + line
			return computeNewYamlAfterPartRemoval(newYaml, baseYamlLines, linePositionToParse+1, fieldPath, 0, ignoreIndentationLevel)
		} else {
			// We are in an interesting identation level

			if strings.HasPrefix(lineTrimmed, (*fieldPath)[currentIdentNumber]+":") {
				if currentIdentNumber == len(*fieldPath)-1 {
					// We found the field or section.
					// Not adding this line to final yaml and setting is deleting this level of identation to true
					return computeNewYamlAfterPartRemoval(newYaml, baseYamlLines, linePositionToParse+1, fieldPath, currentIdentNumber+1, 0)
				} else {
					// We are moving up in the fieldPath but not there yet
					*newYaml = *newYaml + "\n" + line
					return computeNewYamlAfterPartRemoval(newYaml, baseYamlLines, linePositionToParse+1, fieldPath, 0, 0)
				}
			} else {
				// We are in a section not interesting
				// If we are at root, search will continue
				*newYaml = *newYaml + "\n" + line
				return computeNewYamlAfterPartRemoval(newYaml, baseYamlLines, linePositionToParse+1, fieldPath, 0, 0)
			}

		}
	}

}

// Removes a field or section from a yaml file given as a string
// The field must be given with the dot notation.
// Example,
//   - to remove jellyfin section, fieldDotNotation must equal "jellyfin"
//   - to remove jellyfin url field, fieldDotNotation must equal "jellyfin.url"
func RemoveYamlPartHelper(yaml string, fieldDotNotation string) string {
	lines := strings.Split(yaml, "\n")
	fieldPath := strings.Split(fieldDotNotation, ".")
	var newYaml string
	computeNewYamlAfterPartRemoval(&newYaml, &lines, 0, &fieldPath, 0, 0)
	return newYaml
}

func TestLoadContext_ValidConfig(t *testing.T) {
	ctx, err := loadContextFromReader(strings.NewReader(validConfigYAML))

	require.NoError(t, err)
	require.NotNil(t, ctx)
	require.NotNil(t, ctx.Config)
	require.NotNil(t, ctx.Logger)

	assert.Equal(t, "0 8 1 * *", ctx.Config.Scheduler.CronExpr)
	assert.True(t, ctx.Config.Scheduler.Enabled)
	assert.Equal(t, "INFO", ctx.Config.Log.Level)
	assert.Equal(t, "console", ctx.Config.Log.Format)
	assert.Equal(t, "http://localhost:8096", ctx.Config.Jellyfin.Url)
	assert.Equal(t, "secret", ctx.Config.Jellyfin.ApiKey)
	assert.Equal(t, []string{"/movies"}, ctx.Config.Jellyfin.WatchedFilmFolders)
	assert.Equal(t, []string{"/series"}, ctx.Config.Jellyfin.WatchedSeriesFolders)
	assert.Equal(t, 30, ctx.Config.Jellyfin.ObservedPeriodDays)
	assert.Equal(t, false, ctx.Config.Jellyfin.IgnoreItemsAddedAfterLastNewsletter)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30", ctx.Config.Tmdb.ApiKey)
	assert.Equal(t, "smtp.example.com", ctx.Config.SMTP.Host)
	assert.Equal(t, 587, ctx.Config.SMTP.Port)
	assert.Equal(t, "user", ctx.Config.SMTP.Username)
	assert.Equal(t, "pass", ctx.Config.SMTP.Password)
	assert.Equal(t, "Jellyfin", ctx.Config.SMTP.SenderName)
	assert.Equal(t, "classic", ctx.Config.EmailTemplate.Theme)
	assert.Equal(t, "en", ctx.Config.EmailTemplate.Language)
	assert.Equal(t, "New releases", ctx.Config.EmailTemplate.Subject)
	assert.Equal(t, "Newsletter", ctx.Config.EmailTemplate.Title)
	assert.Equal(t, "This week", ctx.Config.EmailTemplate.Subtitle)
	assert.Equal(t, "http://localhost:8096", ctx.Config.EmailTemplate.JellyfinURL)
	assert.Equal(t, "unsub@example.com", ctx.Config.EmailTemplate.UnsubscribeEmail)
	assert.Equal(t, "Admin", ctx.Config.EmailTemplate.JellyfinOwnerName)
	assert.Equal(t, "date_asc", ctx.Config.EmailTemplate.SortMode)
	assert.Equal(t, 10, ctx.Config.EmailTemplate.DisplayOverviewMaxItems)
	assert.Equal(t, true, ctx.Config.DryRun.Enabled)
	assert.Equal(t, false, ctx.Config.DryRun.TestSMTPConnection)
	assert.Equal(t, "/app/config/previews/", ctx.Config.DryRun.OutputDirectory)
	assert.Equal(t, "newsletter_{date}.html", ctx.Config.DryRun.OutputFilename)
	assert.Equal(t, true, ctx.Config.DryRun.IncludeMetadata)
	assert.Equal(t, true, ctx.Config.DryRun.SaveEmailData)
	assert.Equal(t, ctx.Config.EmailRecipients[0], "user1@example.com")
	assert.Equal(t, ctx.Config.EmailRecipients[1], "user2@example.com")
}

func TestLoadContext_MissingRequiredField(t *testing.T) {
	tests := []struct {
		name            string
		yamlKeyToRemove string
		expectedError   string
	}{
		{
			name:            "Missing jellyfin URL",
			yamlKeyToRemove: "jellyfin.url",
			expectedError: `Failed to decode configuration file: [10:9] Key: 'Url' Error:Field validation for 'Url' failed on the 'required' tag
   7 |   level: INFO
   8 |   format: console
   9 | 
> 10 | jellyfin:
               ^
  11 |   api_token: secret
  12 |   watched_film_folders:
  13 |     - /movies
  14 |   `,
		},

		{
			name:            "Missing jellyfin api token",
			yamlKeyToRemove: "jellyfin.api_token",
			expectedError: `Failed to decode configuration file: [10:9] Key: 'ApiToken' Error:Field validation for 'ApiToken' failed on the 'required' tag
   7 |   level: INFO
   8 |   format: console
   9 | 
> 10 | jellyfin:
               ^
  11 |   url: http://localhost:8096
  12 |   watched_film_folders:
  13 |     - /movies
  14 |   `,
		},

		{
			name:            "Missing jellyfin watched_film_folders",
			yamlKeyToRemove: "jellyfin.watched_film_folders",
			expectedError: `Failed to decode configuration file: [10:9] Key: 'WatchedFilmFolders' Error:Field validation for 'WatchedFilmFolders' failed on the 'required' tag
   7 |   level: INFO
   8 |   format: console
   9 | 
> 10 | jellyfin:
               ^
  11 |   url: http://localhost:8096
  12 |   api_token: secret
  13 |   watched_tv_folders:`,
		},

		{
			name:            "Missing jellyfin watched_tv_folders",
			yamlKeyToRemove: "jellyfin.watched_tv_folders",
			expectedError: `Failed to decode configuration file: [10:9] Key: 'WatchedSeriesFolders' Error:Field validation for 'WatchedSeriesFolders' failed on the 'required' tag
   7 |   level: INFO
   8 |   format: console
   9 | 
> 10 | jellyfin:
               ^
  11 |   url: http://localhost:8096
  12 |   api_token: secret
  13 |   watched_film_folders:`,
		},

		{
			name:            "Missing observed_period_days",
			yamlKeyToRemove: "jellyfin.observed_period_days",
			expectedError: `Failed to decode configuration file: [10:9] Key: 'ObservedPeriodDays' Error:Field validation for 'ObservedPeriodDays' failed on the 'required' tag
   7 |   level: INFO
   8 |   format: console
   9 | 
> 10 | jellyfin:
               ^
  11 |   url: http://localhost:8096
  12 |   api_token: secret
  13 |   watched_film_folders:`,
		},

		{
			name:            "Missing TMDB api key",
			yamlKeyToRemove: "tmdb.api_key",
			expectedError:   `Failed to decode configuration file: Key: 'yamlConfiguration.Tmdb.ApiKey' Error:Field validation for 'ApiKey' failed on the 'required' tag`,
		},

		{
			name:            "Missing email_template lang",
			yamlKeyToRemove: "email_template.language",
			expectedError: `Failed to decode configuration file: [31:15] Key: 'Language' Error:Field validation for 'Language' failed on the 'required' tag
  28 |   smtp_sender_email: Jellyfin
  29 |   smtp_tls_type: "TLS"
  30 | 
> 31 | email_template:
                     ^
  32 |   theme: "classic"
  33 |   subject: New releases
  34 |   title: Newsletter
  35 |   `,
		},

		{
			name:            "Missing email_template subject",
			yamlKeyToRemove: "email_template.subject",
			expectedError: `Failed to decode configuration file: [31:15] Key: 'Subject' Error:Field validation for 'Subject' failed on the 'required' tag
  28 |   smtp_sender_email: Jellyfin
  29 |   smtp_tls_type: "TLS"
  30 | 
> 31 | email_template:
                     ^
  32 |   theme: "classic"
  33 |   language: en
  34 |   title: Newsletter
  35 |   `,
		},

		{
			name:            "Missing email_template title",
			yamlKeyToRemove: "email_template.title",
			expectedError: `Failed to decode configuration file: [31:15] Key: 'Title' Error:Field validation for 'Title' failed on the 'required' tag
  28 |   smtp_sender_email: Jellyfin
  29 |   smtp_tls_type: "TLS"
  30 | 
> 31 | email_template:
                     ^
  32 |   theme: "classic"
  33 |   language: en
  34 |   subject: New releases
  35 |   `,
		},

		{
			name:            "Missing email_template subtitle",
			yamlKeyToRemove: "email_template.subtitle",
			expectedError: `Failed to decode configuration file: [31:15] Key: 'Subtitle' Error:Field validation for 'Subtitle' failed on the 'required' tag
  28 |   smtp_sender_email: Jellyfin
  29 |   smtp_tls_type: "TLS"
  30 | 
> 31 | email_template:
                     ^
  32 |   theme: "classic"
  33 |   language: en
  34 |   subject: New releases
  35 |   `,
		},

		{
			name:            "Missing email smtp_server",
			yamlKeyToRemove: "email.smtp_server",
			expectedError: `Failed to decode configuration file: [23:6] Key: 'SmtpServer' Error:Field validation for 'SmtpServer' failed on the 'required' tag
  20 | tmdb:
  21 |   api_key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30
  22 | 
> 23 | email:
            ^
  24 |   smtp_port: 587
  25 |   smtp_username: user
  26 |   smtp_password: pass
  27 |   `,
		},

		{
			name:            "Missing email smtp_port",
			yamlKeyToRemove: "email.smtp_port",
			expectedError: `Failed to decode configuration file: [23:6] Key: 'SmtpPort' Error:Field validation for 'SmtpPort' failed on the 'required' tag
  20 | tmdb:
  21 |   api_key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30
  22 | 
> 23 | email:
            ^
  24 |   smtp_server: smtp.example.com
  25 |   smtp_username: user
  26 |   smtp_password: pass
  27 |   `,
		},

		{
			name:            "Missing email smtp_username",
			yamlKeyToRemove: "email.smtp_username",
			expectedError: `Failed to decode configuration file: [23:6] Key: 'SmtpUsername' Error:Field validation for 'SmtpUsername' failed on the 'required' tag
  20 | tmdb:
  21 |   api_key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30
  22 | 
> 23 | email:
            ^
  24 |   smtp_server: smtp.example.com
  25 |   smtp_port: 587
  26 |   smtp_password: pass
  27 |   `,
		},

		{
			name:            "Missing email smtp_password",
			yamlKeyToRemove: "email.smtp_password",
			expectedError: `Failed to decode configuration file: [23:6] Key: 'SmtpPassword' Error:Field validation for 'SmtpPassword' failed on the 'required' tag
  20 | tmdb:
  21 |   api_key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30
  22 | 
> 23 | email:
            ^
  24 |   smtp_server: smtp.example.com
  25 |   smtp_port: 587
  26 |   smtp_username: user
  27 |   `,
		},

		{
			name:            "Missing email smtp_sender_email",
			yamlKeyToRemove: "email.smtp_sender_email",
			expectedError: `Failed to decode configuration file: [23:6] Key: 'SmtpSenderName' Error:Field validation for 'SmtpSenderName' failed on the 'required' tag
  20 | tmdb:
  21 |   api_key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30
  22 | 
> 23 | email:
            ^
  24 |   smtp_server: smtp.example.com
  25 |   smtp_port: 587
  26 |   smtp_username: user
  27 |   `,
		},

		{
			name:            "Missing recipients",
			yamlKeyToRemove: "recipients",
			expectedError:   `Failed to decode configuration file: Key: 'yamlConfiguration.Recipients' Error:Field validation for 'Recipients' failed on the 'required' tag`,
		},
	}

	finalTests := ""

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			badYamlConfig := RemoveYamlPartHelper(validConfigYAML, tt.yamlKeyToRemove)
			ctx, err := loadContextFromReader(strings.NewReader(badYamlConfig))

			require.NotNil(t, err)

			assert.Nil(t, ctx)
			assert.Equal(t, tt.expectedError, err.Error())
			finalTests = finalTests + "\n{\n name: \"" + tt.name + "\",\nyamlKeyToRemove: \"" + tt.yamlKeyToRemove + "\"\n expectedError: `" + err.Error() + "`,\n},\n"
		})
	}
	//fmt.Println(finalTests)
}
