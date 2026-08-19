# Configuration parameters

Jellyfin Newsletter is configured entirely through a single `config.yml` file. All parameters are **required** unless explicitly marked **optional** below (in the example file, optional parameters are commented out).

Download the example file to get started:

```bash
curl -o config/config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
```

---

## `scheduler`

Controls when the newsletter is sent automatically.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `cron` | Optional | — | Crontab expression defining when the newsletter is sent. Test your expression at [crontab.guru](https://crontab.guru/). |

!> If the entire `scheduler` section is commented out / omitted, the built-in scheduler is disabled and **the newsletter runs once immediately when the container starts**.

```yaml
scheduler:
  cron: "0 8 1 * *" # 8:00 AM on the 1st of every month
```

---

## `jellyfin`

Connection details for your Jellyfin server and which libraries to watch.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `url` | Yes | — | Base URL of your Jellyfin server. |
| `api_token` | Yes | — | Jellyfin API key. See [How to generate a Jellyfin API key](#how-to-generate-a-jellyfin-api-key). |
| `watched_film_folders` | Yes | — | List of movie library folder names to watch for new items. Use **only the last folder name**, without slashes (e.g. `/media/movies` → `movies`). |
| `watched_tv_folders` | Yes | — | List of TV show library folder names to watch for new items. Same naming rule as above (e.g. `/media/tv` → `tv`). |
| `observed_period_days` | Yes | — | Number of days to look back for newly added items. |
| `ignore_item_added_before_last_newsletter` | Optional | `false` | If `true`, items added before the previous newsletter was sent are skipped, even if still inside the observed period. Prevents duplicate mentions across runs. |

```yaml
jellyfin:
  url: ""
  api_token: ""
  watched_film_folders:
    - "movies"
  watched_tv_folders:
    - "tv"
  observed_period_days: 30
  ignore_item_added_before_last_newsletter: false
```

### How to generate a Jellyfin API key
1. Go to your Jellyfin dashboard.
2. Scroll to the **Advanced** section and click **API keys**.
3. Click **+** to create a new API key.
4. Fill in the required fields and save.
5. Copy the generated key into `jellyfin.api_token`.

---

## `tmdb`

| Parameter | Required | Default | Description |
|---|---|---|---|
| `api_key` | Yes | — | TMDB API key, used to fetch posters and movie/show details. See [How to generate a TMDB API key](#how-to-generate-a-tmdb-api-key). |

```yaml
tmdb:
  api_key: ""
```

### How to generate a TMDB API key
1. Go to [themoviedb.org](https://www.themoviedb.org/) and create an account or log in.
2. Open your account settings → **API** section.
3. Click **Create** to generate a new API key.
4. Copy the key named **"API Read Access Token"**.
5. Paste it into `tmdb.api_key`.

---

## `email_template`

Controls the appearance and content of the newsletter email.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `theme` | Yes | — | Theme used to render the email. See [available themes](https://github.com/SeaweedbrainCY/jellyfin-newsletter/tree/main/engine-go/internal/template/themes). Currently: `classic`. |
| `language` | Yes | — | ISO 639 (2-letter) language code for the email content, e.g. `en`, `fr`, `el`. See supported languages on the project's [Weblate](https://weblate.seaweedbrain.xyz). |
| `subject` | Yes | — | Subject line of the email. |
| `title` | Yes | — | Title displayed in the email body. |
| `subtitle` | Yes | — | Subtitle displayed in the email body. |
| `jellyfin_url` | Yes | — | URL used to link back to your Jellyfin instance from the email. |
| `unsubscribe_email` | Yes | — | Contact address shown in the footer's legal/unsubscribe notice. |
| `jellyfin_owner_name` | Yes | — | Name displayed in the email footer. |
| `display_overview_max_items` | Optional | `10` | If the number of new items exceeds this value, item summaries are hidden. `0` = always show summaries, `-1` = always hide summaries. |
| `max_displayed_items` | Optional | disabled (unlimited) | Maximum number of items shown in the email. Extra items are collapsed into a "... and y more" line. Comment out to disable the limit. |
| `sort_mode` | Optional | `date_asc` | Sort order for items in the email. One of `date_asc`, `date_desc`, `name_asc`, `name_desc`. |

Text fields such as `subject`, `title`, and `subtitle` support dynamic placeholders (e.g. `{{.MonthName}}`, `{{.Date}}`) — see the [Placeholders guide](placeholders.md) for the full list and examples.

```yaml
email_template:
  theme: "classic"
  language: "en"
  subject: ""
  title: ""
  subtitle: ""
  jellyfin_url: ""
  unsubscribe_email: ""
  jellyfin_owner_name: ""
  display_overview_max_items: 10
  max_displayed_items: 10
  sort_mode: "date_asc"
```

Want a custom look? See the [Local theme files guide](local-themes.md) to bring your own theme, or [contribute a new one](https://github.com/SeaweedbrainCY/jellyfin-newsletter/blob/main/engine-go/internal/template/themes/README.md) to the built-in set.

---

## `email`

SMTP settings used to send the newsletter. TLS is required.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `smtp_tls_type` | Yes | — | TLS mode: `STARTTLS`, `TLS` (implicit TLS), or `NONE`. |
| `smtp_server` | Yes | — | SMTP server hostname, e.g. `smtp.gmail.com`. |
| `smtp_port` | Yes | — | SMTP port. Typically `587` for STARTTLS or `465` for implicit TLS. |
| `smtp_username` | Yes | — | Username for SMTP authentication. |
| `smtp_password` | Yes | — | Password for SMTP authentication. |
| `smtp_sender_email` | Yes | — | Sender address, e.g. `jellyfin@example.com`, or with a display name: `Jellyfin <jellyfin@example.com>`. |

```yaml
email:
  smtp_tls_type: "STARTTLS"
  smtp_server: ""
  smtp_port:
  smtp_username: ""
  smtp_password: ""
  smtp_sender_email: ""
```

---

## `debug`

| Parameter | Required | Default | Description |
|---|---|---|---|
| `debug` | Optional | `false` | Enables verbose logging and extra diagnostic output. Does not change the script's behavior. |

```yaml
debug: true
```

---

## `log`

| Parameter | Required | Default | Description |
|---|---|---|---|
| `level` | Optional | `INFO` | Minimum log level. One of `DEBUG`, `INFO`, `WARN`, `ERROR`. |
| `format` | Optional | `console` | Log output format. One of `console`, `json`. |

```yaml
log:
  level: INFO
  format: console
```

---

## `dry-run`

Simulates a newsletter run without sending real emails — useful for previewing the rendered HTML or testing your SMTP connection safely. See the [Troubleshooting guide](troubleshooting.md#dry-run-mode) for a full walkthrough.

| Parameter | Required | Default | Description |
|---|---|---|---|
| `enabled` | Optional | `false` | Enables dry-run mode instead of sending real emails. |
| `test_smtp_connection` | Optional | `false` | When `true`, tests the SMTP connection during dry-run without sending an email. When `false`, skips SMTP entirely (preview-only). |
| `output_directory` | Optional | `/app/config/previews/` | Directory where preview files are saved. When using Docker, mount this path to access the generated files. |
| `output_filename` | Optional | `newsletter_{date}.html` | Filename pattern for the generated preview. Supports `{date}`, `{timestamp}`, and `{time}` placeholders. |
| `include_metadata` | Optional | `true` | Includes generation metadata as HTML comments in the preview file. |
| `save_email_data` | Optional | `true` | Also saves a JSON file alongside the HTML preview with all newsletter generation data. |

```yaml
dry-run:
    enabled: false
    test_smtp_connection: false
    output_directory: "/app/config/previews/"
    output_filename: "newsletter_{date}.html"
    include_metadata: true
    save_email_data: true
```

---

## `recipients`

| Parameter | Required | Default | Description |
|---|---|---|---|
| `recipients` | Yes | — | List of email addresses to send the newsletter to. Each entry can be a plain address or `Name <email@example.com>` to set a display name. |

```yaml
recipients:
  - "name@example.com"
  - "Jane Doe <jane@example.com>"
```
