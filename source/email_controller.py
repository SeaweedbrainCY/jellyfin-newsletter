from source import configuration
from source import context
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from source.configuration import logging
from time import sleep
from source.utils import save_last_newsletter_date
import datetime as dt



def send_email(html_content):
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
    
    for recipient in configuration.conf.recipients:
        msg = MIMEMultipart('alternative')
        msg['Subject'] = configuration.conf.email_template.subject.format_map(context.placeholders)
        msg['From'] = configuration.conf.email.smtp_sender_email
        part = MIMEText(html_content, 'html')
    
        msg.attach(part)
        msg['To'] = recipient
        smtp_server.sendmail(configuration.conf.email.smtp_sender_email, recipient, msg.as_string())
        logging.info(f"Email sent to {recipient}")
        sleep(2)
    smtp_server.quit()
    save_last_newsletter_date(dt.datetime.now())

    
