<p align="center">
  <img src="https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/jellyfin_newsletter.png" width="120">
</p>

<h1 align="center">Jellyfin Newsletter</h1>
<p align="center"><strong>Keep your users updated on what's new on your Jellyfin server.</strong></p>

<p align="center">
  <img src="https://github.com/SeaweedbrainCY/jellyfin-newsletter/actions/workflows/build_and_deploy.yml/badge.svg?branch=" />
  <img src="https://img.shields.io/github/license/seaweedbraincy/jellyfin-newsletter" />
  <img src="https://img.shields.io/github/v/release/seaweedbraincy/jellyfin-newsletter" />
</p>

## What is Jellyfin Newsletter?

Jellyfin Newsletter connects to the Jellyfin API to retrieve the movies and TV shows that were recently added to your server, and sends a nicely formatted email digest to your users. It's built in Go, fully self-hosted, and designed to run unattended on a schedule.

<p align="center">
  <img src="https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/new_media_notification_preview.png" width="500">
</p>

## Features

- 🎬 Retrieves the latest added movies and TV shows from your Jellyfin server
- 📬 Sends a newsletter to your configured list of recipients
- 🖼️ Enriches items with details and posters from TMDB
- 📺 Groups TV shows by season
- 🎨 Fully customizable, responsive email template, with support for custom themes
- 🌍 Available in many languages, community-translated on [Weblate](https://weblate.seaweedbrain.xyz)
- 📁 Lets you scope which library folders are watched for new items
- ⏱️ Runs on a schedule via a built-in cron job, or triggered externally

## How it works

1. Jellyfin Newsletter queries your Jellyfin instance for items added within an observation window.
2. Matching items are filtered to the folders you've chosen to watch.
3. Metadata and posters are pulled from TMDB.
4. An HTML email is rendered from the selected theme and sent to your recipients over SMTP.

The whole behavior is controlled by a single `config.yml` file — see the [Configuration parameters](configuration.md) page for the full reference.

## Quick start (Docker)

The recommended way to run Jellyfin Newsletter is with Docker Compose and the built-in cron scheduler.

```bash
curl -o docker-compose.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/docker-compose.yml
mkdir config
curl -o config/config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
```

Edit `config/config.yml` (see [Configuration parameters](configuration.md)), then start the container:

```bash
docker compose up -d
```

?> It's recommended to pin a static image version instead of `latest`, and upgrade manually. See the [releases page](https://github.com/SeaweedbrainCY/jellyfin-newsletter/releases) for the latest version.

You'll also need:

- A **Jellyfin API key** — generate one from your Jellyfin dashboard's Advanced → API keys section
- A **TMDB API key** ("API Read Access Token") — free, from your [TMDB account settings](https://www.themoviedb.org/)
- A working **SMTP server** to send the emails

## Requirements

| Requirement | Notes |
|---|---|
| Docker | Recommended installation method |
| Jellyfin API key | Used to query your Jellyfin instance |
| TMDB API key | Used to fetch posters and metadata |
| SMTP server | TLS required (STARTTLS or implicit TLS) |

## License & contributing

Jellyfin Newsletter is licensed under **AGPLv3**. Contributions are welcome via issues and pull requests — see the project's `CONTRIBUTING.md` for guidelines. If you find it useful, a ⭐️ on [GitHub](https://github.com/SeaweedbrainCY/jellyfin-newsletter) is always appreciated.
