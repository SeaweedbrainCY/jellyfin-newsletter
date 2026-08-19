# Troubleshooting

Jellyfin Newsletter has several built-in tools to help you diagnose setup issues or inspect what a run will produce before it emails anyone.

## Debug mode

Enable verbose logging by setting, in `config.yml`:

```yaml
debug: true
```

This turns on more detailed logging and additional diagnostic output. It does **not** change the script's behavior — it's purely for visibility into what's happening under the hood.

## Dry-run mode

Dry-run mode simulates a newsletter run **without actually sending emails**. Enable it with:

```yaml
dry-run:
    enabled: true
```

By default, this saves preview output to `/app/previews/` for Docker installs (or `./previews/` for local installs):

- The generated HTML email, as `newsletter_{date}_{time}.html`
- Metadata about the configuration and the newsletter generation, embedded in the HTML file
- A JSON metadata file alongside the HTML file with all the information about the generation

Dry-run mode does **not** send any email or test the SMTP connection by default.

?> If you're running via Docker, remember to **mount the `/app/previews/` folder** so you can access the generated preview files from the host.

### Customizing dry-run behavior

All dry-run settings can be tuned directly in `config.yml`:

```yaml
# Dry run settings
dry-run:
    # Enable dry-run mode instead of sending emails
    enabled: false
    # Test SMTP connection during dry-run mode
    # true = dry-run: test SMTP connection but don't send
    # false = preview-only: skip SMTP entirely
    test_smtp_connection: false
    # Directory to save preview files
    output_directory: "/app/config/previews/"
    # Preview filename pattern (supports {date}, {timestamp}, {time} placeholders)
    output_filename: "newsletter_{date}.html"
    # Include generation metadata in HTML comments
    include_metadata: true
    # Save email data as JSON file alongside HTML
    save_email_data: true
```

See the [Configuration parameters — `dry-run`](configuration.md#dry-run) reference for a full breakdown of each field.

## See also

- [Configuration parameters](configuration.md) — full config reference
- [Placeholders](placeholders.md) — dynamic values for your email strings
