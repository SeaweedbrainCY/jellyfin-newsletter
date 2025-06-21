
# Contribution Guidelines for Jellyfin-Newsletter Project

Thank you for your interest in contributing to Jellyfin-Newsletter! We appreciate your willingness to help improve the project. This document outlines the guidelines and processes for contributing to the project, ensuring a smooth and collaborative experience for everyone involved.

## Table of Contents
1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
    - [Fork the Repository](#fork-the-repository)
    - [Set up the Project](#set-up-the-project)
3. [Contribution Process](#contribution-process)
    - [Issue Tracker](#issue-tracker)
    - [Branching](#branching)
    - [Code Style](#code-style)
    - [Testing](#testing)
4. [Submitting Pull Requests](#submitting-pull-requests)
5. [Security](#security)
6. [Community and Communication](#community-and-communication)
7. [Acknowledgements](#acknowledgements)
8. [License](#license)

## Introduction

Jellyfin-Newsletter is an open-source project that aims to provide a comprehensive and user-friendly newsletter system for Jellyfin users. We welcome contributions from developers, designers, and anyone interested in improving the project.

## Getting Started

### Fork the Repository

Before you start contributing, make sure to fork the Jellyfin-Newsletter repository to your GitHub account. You will be working on your forked repository and submitting pull requests from there.

### Requirements
- Python 3.9+ 
- Jellyfin API key - [How to generate an API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-jellyfin-api-key)
- A TMDB API key (free) - [How to generate a TMDB API key](https://github.com/SeaweedbrainCY/jellyfin-newsletter?tab=readme-ov-file#how-to-generate-a-tmdb-api-key)
- A SMTP server 

### Installation and setup

1. Install the required packages 
```bash
# Unix : 
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Windows :
python -m venv venv
venv\Scripts\activate
pip install -r requirements.txt
```
2.  Copy the `config/config-example.yml` to `./config/config.yml` ([direct download](https://raw.githubusercontent.com/SeaweedbrainCY/jellyfin-newsletter/refs/heads/main/config/config-example.yml)) and fill in the required fields. **All fields are required**.
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

3. Make sure to support the following locales:
- `en_US.UTF-8 UTF-8` 
- `fr_FR.UTF-8 UTF-8`

On debian based systems, you can edit `/etc/locale.gen` and uncomment the lines for the locales listed above, then run:
```bash
sudo locale-gen
```
For other systems, please refer to your system's documentation on how to generate locales.

4. Run the script
```bash
python main.py
```

## Contribution Process

### Issue Tracker

If you want to work on a new feature, bug fix, or other enhancements, please check the [Issue Tracker](https://github.com/Jellyfin-Newsletter/Jellyfin-Newsletter/issues) first. It's possible that someone else is already working on something similar or that the issue has already been addressed.

If you find a new issue or want to suggest an enhancement, please open a new issue on the Issue Tracker. Provide a clear description and, if applicable, steps to reproduce the problem.

### Commit signature
We require all commits to be **signed**. [How to](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits). Not-signed commits cannot be merged. To sign your commits, you can use the `-S` flag with `git commit`:
```bash
git commit -S -m "Your commit message"
```

If you already pushed non-signed commits in your forked projects according, you can execute this
```bash 
git rebase --exec 'git commit --amend --no-edit -n -S' -i <first_commit_hash>
```
to sign all commits you've already pushed.

### Branching

For every contribution, create a new branch with a descriptive name that summarizes the changes you plan to make. Use lowercase letters and dashes to separate words, for example:
```
git checkout -b add-localization-support
```


## Submitting Pull Requests

When you're ready to submit your changes, follow these steps:
1. Commit your changes with a descriptive commit message. Commits must be **signed**
2. Push your changes to your forked repository:
```
git push origin your-branch-name
```
3. Go to the Jellyfin-Newsletter repository on GitHub and create a new pull request from your branch.
4. Provide a clear title and description for your pull request, including any relevant information about the changes you made.

## Security

Even for small projects like this one, security is essential. If you discover any security vulnerabilities or potential issues, please **DO NOT** open a public issue. Instead, use the following methods to report security issues:
- Github Private vulnerability reporting (preferred). You can find this option in the "Security" tab of the repository.
- Email at jellynewsletter-security[at]seaweedbrain.xyz. Please encrypt sensitive information using this [PGP public key](https://pgp.stchepinsky.net). 

## Community and Communication

We value a friendly and inclusive community. Be respectful and considerate when communicating with other contributors. If you have questions or need help, through GitHub issues or discussions.

## Acknowledgements

We appreciate all contributions to Jellyfin-Newsletter, and we recognize and acknowledge everyone's effort and time. Contributors will be listed in the project's contributors section.

## License

By contributing to Jellyfin-Newsletter, you agree that your contributions will be licensed under the [GPL-3.0 license](https://github.com/Jellyfin-Newsletter/Jellyfin-Newsletter/blob/main/LICENSE).

---

