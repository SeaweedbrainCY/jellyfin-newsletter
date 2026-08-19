# Run Jellyfin Newsletter with a local theme file

Jellyfin Newsletter ships with [built-in embedded themes](https://github.com/SeaweedbrainCY/jellyfin-newsletter#custom-themes). Before a new theme gets merged into the official set, you can point Jellyfin Newsletter at your own local theme file instead.

## Docker Compose

### Instructions

To use a local theme file:

1. Create the HTML theme file by following the [theme creation guide](https://github.com/SeaweedbrainCY/jellyfin-newsletter/blob/main/engine-go/internal/template/themes/README.md).
2. Store it in a directory (for example `my_themes`) and name the HTML file after your theme (for example `dracula.html` for a theme named `dracula`).
3. Mount that directory into the container (for example, mount `my_themes` at `/app/my_themes`).
4. Update the entrypoint to pass the `-themes-dir` flag pointing to that mount, e.g.:
   ```
   entrypoint: /app/entrypoint -config /app/config/config.yml -themes-dir /app/my_themes
   ```
5. Reference your theme name in `config.yml`, e.g. `theme: "dracula"`.

### Full example

`~/jellyfin-newsletter/my_themes/dracula.html`:

```html
<!doctype html>
<html lang="{{.HTMLLang}}" dir="{{.HTMLDir}}">
    <head>
    </head>
    <body>
         Your super theme
    </body>
</html>
```

`~/jellyfin-newsletter/config/config.yml`:

```yaml
[...]
email_template:
  theme: "dracula" # <-- HERE
  language: "en"
  subject: ""
  title: ""
  subtitle: ""
  jellyfin_url: ""
  unsubscribe_email: ""
  jellyfin_owner_name: ""
[...]
```

`~/jellyfin-newsletter/docker-compose.yml`:

```yaml
services:
  jellyfin-newsletter:
    container_name: jellyfin-newsletter
    image: "ghcr.io/seaweedbraincy/jellyfin-newsletter:1.0.0"
    environment:
      USER_UID: 1000
      USER_GID: 1000
      TZ: "America/Toronto"
    volumes:
      - "./config:/app/config"
      - "./my_themes:/app/my_themes" # <-- HERE
    entrypoint: /app/entrypoint -config /app/config/config.yml -themes-dir /app/my_themes # <-- HERE
```

And you're good to go!

?> **If your local theme directory doesn't contain the theme referenced in your config file, Jellyfin Newsletter falls back to the default embedded themes.** If the theme still can't be found, the script fails to start.

## See also

- [Configuration parameters](configuration.md#email_template) — full `email_template` reference
- [Placeholders guide](placeholders.md) — dynamic values you can use inside your theme and config strings
