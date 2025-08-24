from source import configuration
from source import context
from source.preview_handler import PreviewHandler
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from source.configuration import logging
from time import sleep
from source.utils import save_last_newsletter_date
import datetime as dt


def send_newsletter(html_content, movies=None, series=None, total_tv=0, total_movie=0):
    """
    Send newsletter or generate preview based on configuration
    
    Args:
        html_content (str): Generated HTML email content
        movies (dict): Movies data for metadata
        series (dict): Series data for metadata  
        total_tv (int): Total TV episodes count
        total_movie (int): Total movie count
    
    Returns:
        dict: Result information
    """
    # Check if preview mode is enabled
    if configuration.conf.preview.enabled:
        return _handle_preview_mode(html_content, movies or {}, series or {}, total_tv, total_movie)
    else:
        return _send_normal_email(html_content)


def _handle_preview_mode(html_content, movies, series, total_tv, total_movie):
    """Handle preview or dry-run mode"""
    preview_handler = PreviewHandler()
    
    if configuration.conf.preview.test_smtp_connection:
        # Dry-run mode: test SMTP + save preview
        mode = "dry-run"
        smtp_tested = False
        
        try:
            _test_smtp_connection()
            smtp_tested = True
            logging.info("SMTP connection test: SUCCESS")
        except Exception as e:
            logging.error(f"SMTP connection test: FAILED - {e}")
            
        # Generate metadata with SMTP test result
        metadata = preview_handler.get_metadata(movies, series, total_tv, total_movie, mode, smtp_tested)
        
        # Calculate email size
        metadata['stats']['total_email_size_kb'] = round(len(html_content.encode('utf-8')) / 1024, 1)
        
        # Save preview files
        html_file, json_file = preview_handler.save_preview(html_content, metadata, mode)
        
        # Log dry-run results
        logging.info("DRY-RUN MODE RESULTS:")
        logging.info(f"Would send to: {', '.join(configuration.conf.recipients)}")
        logging.info(f"Email size: {metadata['stats']['total_email_size_kb']}KB")
        logging.info(f"Preview saved: {html_file}")
        if json_file:
            logging.info(f"Metadata saved: {json_file}")
        
        return {
            "mode": "dry-run",
            "smtp_tested": smtp_tested,
            "html_file": html_file,
            "json_file": json_file,
            "recipients": configuration.conf.recipients,
            "email_size_kb": metadata['stats']['total_email_size_kb']
        }
        
    else:
        # Preview-only mode: skip SMTP entirely
        mode = "preview"
        
        # Generate metadata
        metadata = preview_handler.get_metadata(movies, series, total_tv, total_movie, mode, False)
        
        # Calculate email size
        metadata['stats']['total_email_size_kb'] = round(len(html_content.encode('utf-8')) / 1024, 1)
        
        # Save preview files
        html_file, json_file = preview_handler.save_preview(html_content, metadata, mode)
        
        # Log preview results
        logging.info("PREVIEW MODE RESULTS:")
        logging.info(f"Email size: {metadata['stats']['total_email_size_kb']}KB")
        logging.info(f"Preview saved: {html_file}")
        if json_file:
            logging.info(f"Metadata saved: {json_file}")
        logging.info("SMTP testing skipped (preview-only mode)")
        
        return {
            "mode": "preview",
            "smtp_tested": False,
            "html_file": html_file,
            "json_file": json_file,
            "email_size_kb": metadata['stats']['total_email_size_kb']
        }


def _test_smtp_connection():
    """Test SMTP connection without sending email"""
    tls_type = configuration.conf.email.smtp_tls_type.upper()
    
    if tls_type == "TLS":
        smtp_server = smtplib.SMTP_SSL(configuration.conf.email.smtp_server, configuration.conf.email.smtp_port)
    elif tls_type == "STARTTLS":
        smtp_server = smtplib.SMTP(configuration.conf.email.smtp_server, configuration.conf.email.smtp_port)
        smtp_server.starttls()
    else:
        raise Exception(f"Invalid SMTP TLS type: {tls_type}")
    
    # Test login
    smtp_server.login(configuration.conf.email.smtp_user, configuration.conf.email.smtp_password)
    
    # Test recipient validation (basic check)
    for recipient in configuration.conf.recipients:
        if '@' not in recipient:
            raise Exception(f"Invalid recipient email: {recipient}")
    
    smtp_server.quit()
    logging.info(f"SMTP server: {configuration.conf.email.smtp_server}:{configuration.conf.email.smtp_port}")
    logging.info(f"Recipients validated: {len(configuration.conf.recipients)} addresses")


def _send_normal_email(html_content):
    """Send email normally (original functionality)"""
    try:      
        tls_type = configuration.conf.email.smtp_tls_type.upper()
        if tls_type == "TLS":
            smtp_server = smtplib.SMTP_SSL(configuration.conf.email.smtp_server, configuration.conf.email.smtp_port)
        elif tls_type == "STARTTLS":
            smtp_server = smtplib.SMTP(configuration.conf.email.smtp_server, configuration.conf.email.smtp_port)
            smtp_server.starttls()
        else:
            raise Exception(f"Invalid SMTP TLS type: {tls_type}")
        smtp_server.login(configuration.conf.email.smtp_user, configuration.conf.email.smtp_password)
    except Exception as e:
        raise Exception(f"Error while connecting to the SMTP server. Got error: {e}")
    
    sent_count = 0
    for recipient in configuration.conf.recipients:
        msg = MIMEMultipart('alternative')
        msg['Subject'] = configuration.conf.email_template.subject.format_map(context.placeholders)
        msg['From'] = configuration.conf.email.smtp_sender_email
        part = MIMEText(html_content, 'html')
    
        msg.attach(part)
        msg['To'] = recipient
        smtp_server.sendmail(configuration.conf.email.smtp_sender_email, recipient, msg.as_string())
        logging.info(f"Email sent to {recipient}")
        sent_count += 1
        sleep(2)
    smtp_server.quit()
    save_last_newsletter_date(dt.datetime.now())
    
    return {
        "mode": "normal",
        "sent_count": sent_count,
        "recipients": configuration.conf.recipients
    }


# Legacy function for backwards compatibility
def send_email(html_content):
    """Legacy function - use send_newsletter instead"""
    result = _send_normal_email(html_content)
    return result
