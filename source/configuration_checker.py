from source.configuration import conf
from source.configuration import logging
import re
from urllib.parse import urlparse

def check_jellyfin_configuration():
    
    # Jellyfin URL
    parsed_url = urlparse(conf.jellyfin.url)
    assert parsed_url.scheme != '', f"[FATAL] Invalid Jellyfin URL. The URL must contain the scheme (e.g. http:// or https://). Please check the configuration. Got {conf.jellyfin.url}. Parsed : {parsed_url}"
    assert parsed_url.netloc != '', f"[FATAL] Invalid Jellyfin URL. The URL must contain a valid host (e.g. example.com or 127.0.0.1:80). Please check the configuration. Got {conf.jellyfin.url}. Parsed : {parsed_url}"

    # Jellyfin API token
    assert isinstance(conf.jellyfin.api_token, str), "[FATAL] Invalid Jellyfin API token. The API token must be a string. Please check the configuration."
    assert conf.jellyfin.api_token != '', "[FATAL] Invalid Jellyfin API token. The API token cannot be empty. Please check the configuration."

    # watched_film_folders 
    assert isinstance(conf.jellyfin.watched_film_folders, list), "[FATAL] Invalid watched film folders. The watched film folders must be a list. Please check the configuration."
    
    # watched_tv_folders
    assert isinstance(conf.jellyfin.watched_tv_folders, list), "[FATAL] Invalid watched TV folders. The watched TV folders must be a list. Please check the configuration."

    # observed_period_days
    assert isinstance(conf.jellyfin.observed_period_days, int), "[FATAL] Invalid observed period days. The observed period days must be an integer. Please check the configuration."

    # ignore_item_added_before_last_newsletter
    assert isinstance(conf.jellyfin.ignore_item_added_before_last_newsletter, bool), "[FATAL] Invalid ignore_item_added_before_last_newsletter. The ignore_item_added_before_last_newsletter must be a boolean. Please check the configuration."



def check_tmdb_configuration():
    # TMDB API key
    assert isinstance(conf.tmdb.api_key, str), "[FATAL] Invalid TMDB API key. The API key must be a string. Please check the configuration."
    assert conf.tmdb.api_key != '', "[FATAL] Invalid TMDB API key. The API key cannot be empty. Please check the configuration."


def email_template_configuration():
    # Language
    assert isinstance(conf.email_template.language, str), "[FATAL] Invalid email template language. The language must be a string. Please check the configuration."
    assert conf.email_template.language in ['en', 'fr'], "[FATAL] Invalid email template language. The language must be either 'en' or 'fr'. Please check the configuration."

    # Subject
    assert isinstance(conf.email_template.subject, str), "[FATAL] Invalid email template subject. The subject must be a string. Please check the configuration."

    # Title
    assert isinstance(conf.email_template.title, str), "[FATAL] Invalid email template title. The title must be a string. Please check the configuration."

    # Subtitle
    assert isinstance(conf.email_template.subtitle, str), "[FATAL] Invalid email template subtitle. The subtitle must be a string. Please check the configuration."

    # Jellyfin URL
    assert isinstance(conf.email_template.jellyfin_url, str), "[FATAL] Invalid email template Jellyfin URL. The Jellyfin URL must be a string. Please check the configuration."

    # Unsubscribe email
    assert isinstance(conf.email_template.unsubscribe_email, str), "[FATAL] Invalid email template unsubscribe email. The unsubscribe email must be a string. Please check the configuration."

    # Jellyfin owner name
    assert isinstance(conf.email_template.jellyfin_owner_name, str), "[FATAL] Invalid email template Jellyfin owner name. The Jellyfin owner name must be a string. Please check the configuration."


def check_email_configuration():
    # SMTP server
    assert isinstance(conf.email.smtp_server, str), "[FATAL] Invalid email SMTP server. The SMTP server must be a string. Please check the configuration."
    assert conf.email.smtp_server != '', "[FATAL] Invalid email SMTP server. The SMTP server cannot be empty. Please check the configuration."

    # SMTP port
    assert isinstance(conf.email.smtp_port, int), "[FATAL] Invalid email SMTP port. The SMTP port must be an integer. Please check the configuration."
    assert conf.email.smtp_port > 0, "[FATAL] Invalid email SMTP port. The SMTP port must be greater than 0. Please check the configuration."

    # SMTP username
    assert isinstance(conf.email.smtp_user, str), "[FATAL] Invalid email SMTP username. The SMTP username must be a string. Please check the configuration."
    assert conf.email.smtp_user != '', "[FATAL] Invalid email SMTP username. The SMTP username cannot be empty. Please check the configuration."

    # SMTP password
    assert isinstance(conf.email.smtp_password, str), "[FATAL] Invalid email SMTP password. The SMTP password must be a string. Please check the configuration."
    assert conf.email.smtp_password != '', "[FATAL] Invalid email SMTP password. The SMTP password cannot be empty. Please check the configuration."
    
    # SMTP TLS type
    assert isinstance(conf.email.smtp_tls_type, str), "[FATAL] Invalid email SMTP TLS type. The SMTP TLS type must be a string. Please check the configuration."
    assert conf.email.smtp_tls_type in ['STARTTLS', 'TLS'], "[FATAL] Invalid SMTP TLS type. The SMTP TLS type must be either 'STARTTLS' or 'TLS'. Please check the configuration."
    
def check_recipients_configuration():
    # Recipients
    assert isinstance(conf.recipients, list), "[FATAL] Invalid recipients configuration. The recipients must be a list. Please check the configuration."
    

def check_scheduler_configuration():
    if conf.scheduler.enabled:
        assert isinstance(conf.scheduler.cron, str), "[FATAL] Invalid scheduler cron expression. The cron expression must be a string. Please check the configuration."


def check_preview_configuration():
    # enabled
    assert isinstance(conf.preview.enabled, bool), "[FATAL] Invalid preview.enabled. The enabled flag must be a boolean. Please check the configuration."
    
    # test_smtp_connection  
    assert isinstance(conf.preview.test_smtp_connection, bool), "[FATAL] Invalid preview.test_smtp_connection. The test_smtp_connection flag must be a boolean. Please check the configuration."
    
    # output_directory
    assert isinstance(conf.preview.output_directory, str), "[FATAL] Invalid preview.output_directory. The output_directory must be a string. Please check the configuration."
    assert conf.preview.output_directory != '', "[FATAL] Invalid preview.output_directory. The output_directory cannot be empty. Please check the configuration."
    
    # output_filename
    assert isinstance(conf.preview.output_filename, str), "[FATAL] Invalid preview.output_filename. The output_filename must be a string. Please check the configuration."
    assert conf.preview.output_filename != '', "[FATAL] Invalid preview.output_filename. The output_filename cannot be empty. Please check the configuration."
    
    # include_metadata
    assert isinstance(conf.preview.include_metadata, bool), "[FATAL] Invalid preview.include_metadata. The include_metadata flag must be a boolean. Please check the configuration."
    
    # save_email_data
    assert isinstance(conf.preview.save_email_data, bool), "[FATAL] Invalid preview.save_email_data. The save_email_data flag must be a boolean. Please check the configuration."
    
    # If preview enabled, validate directory exists or can be created
    if conf.preview.enabled:
        import os
        try:
            os.makedirs(conf.preview.output_directory, exist_ok=True)
        except Exception as e:
            logging.error(f"[FATAL] Cannot create preview output directory '{conf.preview.output_directory}': {e}")
            exit(1)


def check_configuration():
    """
    Check if the configuration is valid.
    The goal is to ensure all values fetched from the configuration file are valid.
    """
    check_jellyfin_configuration()
    check_tmdb_configuration()
    email_template_configuration()
    check_email_configuration()
    check_recipients_configuration()
    check_preview_configuration()
    
    