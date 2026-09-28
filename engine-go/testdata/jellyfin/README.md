## Jellyfin test data

This directory contains the Docker setup and mounted volumes used to run a **test-only Jellyfin instance** for generating `go-vcr` cassettes.

Use it only when fixtures need to be refreshed (for example after upgrading Jellyfin, changing API usage, or updating expected newsletter content).

### Notes

- This environment is test-specific and uses fake media files.
- API keys/tokens/config values are temporary and test-only.
- Media files are created at runtime so Jellyfin `DateCreated` values are deterministic for integration tests.
- If fixture data changes, update test assertions accordingly.

---

## Regenerate cassettes (integration fixtures)

### 1) Start Jellyfin test container

From this directory (`engine-go/testdata/jellyfin`), start the stack:

```sh
docker-compose up -d
```

### 2) Configure Jellyfin manually

Open Jellyfin and complete initial setup.  
Create exactly 2 libraries:

- Movies library -> `/movies`
- TV Shows library -> `/tvshows`

### 3) Create API key and update integration config

Create a Jellyfin API key, then temporarily update:

`engine-go/testdata/fixtures/config/config.test.yml`

with the correct Jellyfin URL and API key.

### 4) Run integration test from `engine-go`

```sh
cd engine-go
export INTEGRATION_TEST_CONFIG_FILE=../../testdata/fixtures/config/config.test.yml
go test -tags=integration ./internal/newsletter -run TestJellyfinNewsletter -v
```

This will record/update cassettes under `engine-go/testdata/fixtures/`.

### 5) Update fixed test date

Update the fixed date in:

`engine-go/internal/newsletter/newsletter_integration_test.go`

Specifically, adjust `FakeClock.Now()` and expected addition date to the date used for the newly recorded fixtures.

---

## After recording

- Re-run the same test once more to ensure playback is stable.
- Review cassette diffs and update assertions in `newsletter_integration_test.go` if expected content changed.
- Stop containers when finished:

```sh
docker-compose down
```
- Revert the config.test.yml
