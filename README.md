# Jellyfin Newsletter - keep your users updated

<p align="center">
<img src="https://github.com/SeaweedbrainCY/jellyfin-newsletter/actions/workflows/build_and_deploy.yml/badge.svg?branch="/>
 <img src="https://img.shields.io/github/license/seaweedbraincy/jellyfin-newsletter"/>
<img src="https://img.shields.io/github/v/release/seaweedbraincy/jellyfin-newsletter"/>
</p>

<p align="center">
<img src="https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/assets/jellyfin_newsletter.png" width=100>
</p>

A newsletter for Jellyfin to notify your users of your latest additions. Jellyfin Newsletter connects to the Jellyfin API to retrieve recently added items and send them to your users. 

It is fully customizable and can be run on a schedule using a cron job or a task scheduler.

## Table of Contents
1. [What it looks like](#what-it-looks-like)
2. [Features](#features)
3. [Recommended installation: Docker](#recommended-installation-docker)
4. [Current limitations](#current-limitations)
5. [License](#license)
6. [Contribution](#contribution)
7. [How to](#how-to)
   - [How to generate a Jellyfin API key](#how-to-generate-a-jellyfin-api-key)
   - [How to generate a TMDB API key](#how-to-generate-a-tmdb-api-key)

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
- Support for English and French
- Configure the list of recipients
- Configure specific folders to watch for new items

## Recommended installation: Docker
### Requirements

- Docker 
- Jellyfin API key - [How to generate an API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-jellyfin-api-key)
- A TMDB API key (free) - [How to generate a TMDB API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-tmdb-api-key)
- A SMTP server 

### Configuration
1. Download the [config file](https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml):
```
curl -o config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
```
2. Edit the `config.yml` file and fill in the required fields. **All fields are required.**
```yaml
jellyfin:
    # URL of your jellyfin server
    url: "" 

    # API token of your jellyfin server, see requirements for more info
    api_token: ""

    # List of folders to watch for new movies 
    # You can find them in your Jellyfin Dashboard -> Libraries -> Select a library -> Folder **WITHOUT THE TRAILING /**
    watched_film_folders:
        - ""
        # example for /movies folder add "movies"


    # List of folders to watch for new shows
    # You can find them in your Jellyfin Dashboard -> Libraries -> Select a library -> Folder **WITHOUT THE TRAILING /**
    watched_tv_folders:
        - ""
        # example for /tv folder add "tv"
  
  # Number of days to look back for new items
  observed_period_days: 30

tmdb:
    # TMDB API key, see requirements for more info
    api_key: ""

# Email template to use for the newsletter
# You can use placeholders to dynamically insert values. See available placeholders here : https://github.com/SeaweedbrainCY/jellyfin-newsletter/wiki/How-to-use-placeholder
email_template:
    # Language of the email, supported languages are "en" and "fr"
    language: "en"
    # Subject of the email
    subject: ""
    # Title of the email
    title: ""
    # Subtitle of the email
    subtitle: ""
    # Will be used to redirect the user to your Jellyfin instance
    jellyfin_url: ""
    # For the legal notice in the footer
    unsubscribe_email: ""
    # Used in the footer
    jellyfin_owner_name: ""

# SMTP server configuration, TLS is required for now
# Check your email provider for more information
email:
    # Example: GMail: smtp.gmail.com
    smtp_server: ""
    # Usually 587
    smtp_port: 
    # The username of your SMTP account
    smtp_username: ""
    # The password of your SMTP account
    smtp_password: ""
    # Example: "jellyfin@example.com" or to set display username "Jellyfin <jellyfin@example.com>"
    smtp_sender_email: ""


# List of users to send the newsletter to
recipients:
  - ""
  # Example: "name@example.com" or to set username "Name <name@example.com>"
```

3. Run the docker container 
```bash
docker run --rm \
    -v ./config.yml:/app/config/config.yml \
    ghcr.io/seaweedbraincy/jellyfin-newsletter:latest
```
*Note: It is recommended to use a static version instead of  `latest`, and manually upgrade. [Last version](https://github.com/SeaweedbrainCY/jellyfin-newsletter/pkgs/container/jellyfin-newsletter)*

4. (Optional) Schedule the script to run on a regular basis. 
```bash
# Unix :
crontab -e
# Add the following line to run the script every 1st of the month at 8am
0 8 1 * * root docker run --rm -v $(pwd)/config.yml:/app/config/config.yml ghcr.io/seaweedbraincy/jellyfin-newsletter:latest
```
*Note: It is recommended to use a static version instead of `latest`, and manually upgrade. [Last version](https://github.com/SeaweedbrainCY/jellyfin-newsletter/pkgs/container/jellyfin-newsletter)*




## Current limitations
- Only supports English and French languages for the email template
- Only supports TLS for the SMTP server
- Only supports movies and TV shows for now
- Not available as a Jellyfin plugin yet 
- Must be run manually or scheduled

## License
This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

## Contribution
Feel free to contribute to this project by opening an issue or a pull request.

A contribution guide is available in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

If you like this project, consider giving it a ⭐️.

If you encounter any issues, please let me know by opening an issue.

## How to 
### How to generate a Jellyfin API key
1. Go to your Jellyfin dashboard
2. Scroll to advanced section and click on API keys
3. Click on the `+` button to create a new API key
4. Fill in the required fields and click on save
5. Copy the generated API key
6. Paste it in the `config.yml` file under `jellyfin.api_token`

### How to generate a TMDB API key
1. Go to the [TMDB website](https://www.themoviedb.org/)
2. Create an account or log in
3. Go to your account settings
4. Click on the API section
5. Click on the `Create` button to create a new API key
6. Copy the API key named "API Read Access Token"
7. Paste it in the `config.yml` file under `tmdb.api_key`
