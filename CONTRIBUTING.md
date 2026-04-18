# Contribution Guidelines for Jellyfin-Newsletter

Thank you for your interest in contributing to Jellyfin-Newsletter! This document reflects the current Go-based engine (`engine-go`).

> **Note on the Python codebase:** The original Python implementation is kept in the repository for historical reference only. It is deprecated and will be removed in a future release. All active development happens in `engine-go`.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Contribution Process](#contribution-process)
3. [Code Style & Linting](#code-style--linting)
4. [Testing](#testing)
5. [Translations](#translations)
6. [Submitting Pull Requests](#submitting-pull-requests)
7. [Security](#security)
8. [License](#license)

---

## Getting Started

### Requirements

- Go 1.23+
- A running Jellyfin instance with an API key — [How to generate a Jellyfin API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-jellyfin-api-key)
- A TMDB API key (free) — [How to generate a TMDB API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-tmdb-api-key)
- An SMTP server
- Docker (required for integration tests)

### Setup

1. Fork and clone the repository.

2. Navigate to the Go engine directory:
   ```bash
   cd engine-go
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Copy the example config and fill in the required fields:
   ```bash
   cp config/config-example.yml config/config.yml
   ```

5. Run the application:
   ```bash
   go run . --config config/config.yml
   ```

---

## Contribution Process

### Issue Tracker

Before starting work, check the [Issue Tracker](https://github.com/SeaweedbrainCY/jellyfin-newsletter/issues). If no existing issue covers your change, open one with a clear description and — for bugs — steps to reproduce.

### Branching

Create a descriptive branch for your work, using lowercase letters and dashes:

```bash
git checkout -b add-localization-support
```

### Commit Signatures

All commits must be **signed**. See [GitHub's signing guide](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).

```bash
git commit -S -m "your commit message"
```

If you already pushed unsigned commits, you can sign them retroactively:

```bash
git rebase --exec 'git commit --amend --no-edit -n -S' -i <first_commit_hash>
```

---

## Code Style & Linting

The project uses [golangci-lint](https://golangci-lint.run/). Before submitting, make sure your code passes:

```bash
# from engine-go/
make lint
```

Fix any reported issues before opening a PR. PRs with lint failures will not be merged.
Linter exceptions should be avoided as much as possible. If it is absolutely necessary, it needs to be discussed with mantainers to be accepted.

To format code, you can use 
```bash
# from engine-go/
make fmt
```

---

## Testing

The project has both unit tests and integration tests. Integration tests use Docker (via testcontainers) and are tagged with `//go:build integration`.

```bash
# from engine-go/

# Unit tests only
make test

# Integration tests (requires Docker)
make integration
```

Contributions that introduce new features or fix bugs should include appropriate test coverage.

---

## Translations

Translations are managed via **Weblate** at [weblate.seaweedbrain.xyz](https://weblate.seaweedbrain.xyz). If you want to add or improve a translation:

1. Head to the Weblate project and contribute there directly — no PR needed for translation-only changes.

Do **not** edit translation files manually in the repository — they are synced from Weblate.

---

## Submitting Pull Requests

1. Make sure `make lint` and `make test` both pass locally.
2. Push your branch to your fork:
   ```bash
   git push origin your-branch-name
   ```
3. Open a pull request against the main repository with a clear title and description of your changes.
4. Reference any related issues in the PR description.

---

## Security

Do **not** open a public issue for security vulnerabilities. Instead:

- **Preferred:** Use [GitHub Private Vulnerability Reporting](https://github.com/SeaweedbrainCY/jellyfin-newsletter/security) (Security tab of the repository).
- **Alternative:** Email `jellynewsletter-security[at]seaweedbrain.xyz`. For sensitive information, please encrypt using the [PGP public key](https://pgp.stchepinsky.net).

---

## Community and Communication

Be respectful and considerate. Questions and discussions can happen through GitHub Issues or Discussions.

---

## License

Jellyfin-Newsletter is licensed under **AGPLv3**. By contributing, you agree that your contributions will be licensed under the same terms. See the [LICENSE](LICENSE) file for details.
