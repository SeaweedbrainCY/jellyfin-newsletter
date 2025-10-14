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
- Support for English, French, and Hebrew (RTL-aware template)
- Configure the list of recipients
- Configure specific folders to watch for new items

## Recommended installation: Docker
### Requirements

- Docker 
- Jellyfin API key - [How to generate an API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-jellyfin-api-key)
- A TMDB API key (free) - [How to generate a TMDB API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-tmdb-api-key)
- A SMTP server 

### Configuration with built-in cron job
This is the default and recommended way to run the newsletter. The Docker container will run on a schedule using a built-in cron job. It will run on the schedule defined in the `config/config.yml` file.

1. Download the [docker-compose.yml](https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/docker-compose.yml) file:
```bash 
curl -o docker-compose.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/docker-compose.yml
```

2. (optional) Edit the `docker-compose.yml` file to change the default user or timezone.

3. Create a `config` folder in the same directory as the `docker-compose.yml` file:
```bash
mkdir config
```

4. Download the [config file](https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml) in the `config` folder:
```
curl -o config/config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
```

5. Edit the `config/config.yml` file and fill in the required fields. **All non-commented fields are required.**

6. Run the docker container with docker compose 
```bash
docker compose up -d
```

> [!note]
> Note: It is recommended to use a static version instead of `latest`, and manually upgrade. [Last version](https://github.com/SeaweedbrainCY/jellyfin-newsletter/releases)


### Configuration with external cron job
Use this method if you want to run the script on a schedule using an external cron job or task scheduler, instead of the built-in cron job. Docker will run once, and exit after sending the newsletter.

1. Create a `config` folder.
```bash
mkdir config
```

2. Download the [config file](https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml) in the `config` folder:
```
curl -o config/config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
```

2. Edit the `config/config.yml` file and fill in the required fields. **All non-commented fields are required.**


3. Run the docker container to send the newsletter
```bash
docker run --rm \
    -v ./config:/app/config \
    ghcr.io/seaweedbraincy/jellyfin-newsletter:v0.9.0
```
> [!note]
> Note: It is recommended to use a static version instead of `latest`, and manually upgrade. [Last version](https://github.com/SeaweedbrainCY/jellyfin-newsletter/releases)

4. Schedule the script to run on a regular basis. 
```bash
# Unix :
crontab -e
# Add the following line to run the script every 1st of the month at 8am
0 8 1 * * root docker run --rm -v PATH_TO_CONFIG_FOLDER/config:/app/config/ ghcr.io/seaweedbraincy/jellyfin-newsletter:v0.9.0
```





## Current limitations
- Email template translations are available for English, French, and Hebrew
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
