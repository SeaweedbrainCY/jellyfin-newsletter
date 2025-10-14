"""
Add context management for placeholders. 
Here are defined all placeholders the user can use in custom string to customize their email. 
"""


import datetime as dt 
from source import configuration
import locale

class SafeFormatDict(dict):
    """
    A dictionary that allows safe formatting of strings with placeholders.
    If a key is not found, it returns a placeholder string instead of raising a KeyError.
    """
    def __missing__(self, key): 
        return key.join("{}")

# Set locale to the user's locale
if configuration.conf.email_template.language == "fr":
    locale.setlocale(locale.LC_TIME, 'fr_FR.UTF-8')
elif configuration.conf.email_template.language == "he":
    locale.setlocale(locale.LC_TIME, 'he_IL.UTF-8')
else:
    locale.setlocale(locale.LC_TIME, 'en_US.UTF-8')

placeholders = SafeFormatDict({
    "date": dt.datetime.now().strftime("%Y-%m-%d"),
    "day_name": dt.datetime.now().strftime("%A"),
    "day_number": dt.datetime.now().strftime("%d"),
    "month_name": dt.datetime.now().strftime("%B"),
    "month_number": dt.datetime.now().strftime("%m"),
    "year": dt.datetime.now().strftime("%Y"),
    "start_date": (dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)).strftime("%Y-%m-%d"),
    "start_day_name": (dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)).strftime("%A"),
    "start_day_number": (dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)).strftime("%d"),
    "start_month_name": (dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)).strftime("%B"),
    "start_month_number": (dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)).strftime("%m"),
    "start_year": (dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)).strftime("%Y")

})