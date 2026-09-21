# Jellyfin Newsletter - keep your users updated

<p align="center">
    <img src="https://img.shields.io/github/stars/seaweedbraincy/jellyfin-newsletter" alt="GitHub stars" style="margin-right:10px" /> 
    <img src="https://img.shields.io/badge/GHCR%20pulls-48.9K-blue?logo=docker" alt="GHCR pulls" />
    <img src="https://img.shields.io/github/license/seaweedbraincy/jellyfin-newsletter"/>
<img src="https://img.shields.io/github/v/release/seaweedbraincy/jellyfin-newsletter"/>
</p>

<p align="center">
<img src="https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/jellyfin_newsletter.png" width=100>
</p>

A newsletter for Jellyfin to notify your users of your latest additions. Jellyfin Newsletter connects to the Jellyfin API to retrieve recently added items and send them to your users. 

It is fully customizable and can be run on a schedule using a cron job or a task scheduler.

> [!warning]
> **Upgrading to Jellyfin 12 ?**
>
> [Read this first](https://github.com/SeaweedbrainCY/jellyfin-newsletter/discussions/162)

## Table of Contents
1. [What it looks like](#what-it-looks-like)
2. [Features](#features)
3. [Getting started](#getting-started)
4. [License](#license)
5. [Contribution](#contribution)

## What it looks like 
<p align="center">
<img src="https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/new_media_notification_preview.png" width=500>
</p>

## Features
- Retrieve the last added movies and TV shows from your Jellyfin server
- Send a newsletter to your users with the last added items
- Retrieve the movie details from TMDB, including poster
- Group TV shows by seasons
- Fully customizable and responsive email template
- Easy to maintain, extend, setup and run
- Support many languages (see below)
- Configure the list of recipients
- Configure specific folders to watch for new items
- Support themes 

### Supported languages
You can contribute to the translation of Jellyfin-Newsletter on [Crowdin](https://crowdin.com/project/jellyfin-Newsletter)
<p align="center">
<a href="https://crowdin.com/project/jellyfin-Newsletter">
<img src="https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/crowdin-status.svg" alt="State of localization" />
</a>
</p>

### Custom themes
#### Create a new theme
You can create and propose a new theme by following the [theme creation guide](engine-go/internal/template/themes/README.md).

Currently available themes:
- `Classic`

#### Bring your own themes
Jellyfin-Newsletter comes with built-in and embedded themes (listed above). However, before your new beautiful theme gets included in the official themes, you can just instruct Jellyfin-Newsletter to use a local theme file. 

[Follow these instructions to use your own local theme files](https://github.com/SeaweedbrainCY/jellyfin-newsletter/wiki/Bring-your-own-themes)


## Getting started

Follow the [official documentation](https://jellyfin-newsletter.seaweedbrain.xyz/) to get started.


## License
This project is licensed under the AGPLv3 License—see the [LICENSE](LICENSE) file for details.

## Contribution
Feel free to contribute to this project by opening an issue or a pull request.

A contribution guide is available in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

If you like this project, consider giving it a ⭐️.

If you encounter any issues, please let me know by opening an issue.
