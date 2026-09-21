# Installation

Jellyfin Newsletter ships as a Docker image.

## Docker 

### Requirements

- Docker
- A Jellyfin API key — [how to generate one](./configuration#how-to-generate-a-jellyfin-api-key)
- A TMDB API key (free) — [how to generate one](./configuration#how-to-generate-a-tmdb-api-key)
- An SMTP server to send from

There are two ways to run it, depending on how you want scheduling handled.

| | Container stays running | Schedule defined in | Recommended for|
| --- | --- | --- | --- |
| Using built-in cron | Yes  |`config.yml`| Most setups | 
| Using external cron | No — runs once, sends, exits |Your own crontab / task scheduler | Setups that already centralize cron/systemd timers|



### Option A: Built-in cron job (recommended)

The container stays up and sends on the schedule defined in `config/config.yml`.

1. Download the `docker-compose.yml` file:

   ```bash
   curl -o docker-compose.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/docker-compose.yml
   ```

2. *(Optional)* Edit `docker-compose.yml` to change the default user or timezone.

3. Create a `config` folder next to it:

   ```bash
   mkdir config
   ```

4. Download the example config into that folder:

   ```bash
   curl -o config/config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
   ```

5. Edit `config/config.yml` and fill in the required fields. **Every non-commented field is required.** See the full [Configuration reference](./configuration.md).

6. Start it:

   ```bash
   docker compose up -d
   ```

::: tip
Pin a static version instead of `latest`, and upgrade manually. See the [releases page](https://github.com/SeaweedbrainCY/jellyfin-newsletter/releases) for the latest tag.
:::

### Option B: External cron job

Use this if you'd rather trigger the send yourself — the container runs once and exits after sending.

1. Create a `config` folder:

   ```bash
   mkdir config
   ```

2. Download the example config into it:

   ```bash
   curl -o config/config.yml https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml
   ```

3. Edit `config/config.yml` and fill in the required fields. **Every non-commented field is required.**

4. Run it once, to send immediately:

   ```bash
   docker run --rm \
       -v ./config:/app/config \
       ghcr.io/seaweedbraincy/jellyfin-newsletter:1.3.0
   ```

5. Schedule it to run regularly. For example, with `cron`, to send on the 1st of every month at 8am:

   ```bash
   crontab -e
   # then add:
   0 8 1 * * root docker run --rm -v PATH_TO_CONFIG_FOLDER/config:/app/config/ ghcr.io/seaweedbraincy/jellyfin-newsletter:1.3.0
   ```

::: warning
Replace `PATH_TO_CONFIG_FOLDER` with the absolute path to your `config` folder, and pin a real version tag instead of `latest`.
:::

## Next steps

- [Configuration](./configuration.md) — every field in `config.yml`, explained
