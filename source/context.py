"""
Add context management for placeholders. 
Here are defined all placeholders the user can use in custom string to customize their email. 
"""


import datetime as dt 
from source import configuration
from source.configuration import logging

class SafeFormatDict(dict):
    """
    A dictionary that allows safe formatting of strings with placeholders.
    If a key is not found, it returns a placeholder string instead of raising a KeyError.
    """
    def __missing__(self, key): 
        return key.join("{}")

# Localized day/month names without relying on system locales
_DAY_NAMES = {
    "en": ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"],
    "fr": ["lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi", "dimanche"],
    "he": ["יום שני", "יום שלישי", "יום רביעי", "יום חמישי", "יום שישי", "שבת", "יום ראשון"],
}

_MONTH_NAMES = {
    "en": [
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"
    ],
    "fr": [
        "janvier", "février", "mars", "avril", "mai", "juin",
        "juillet", "août", "septembre", "octobre", "novembre", "décembre"
    ],
    "he": [
        "ינואר", "פברואר", "מרץ", "אפריל", "מאי", "יוני",
        "יולי", "אוגוסט", "ספטמבר", "אוקטובר", "נובמבר", "דצמבר"
    ],
}

def _day_name(date: dt.datetime, lang: str) -> str:
    names = _DAY_NAMES.get(lang) or _DAY_NAMES.get("en")
    # Python weekday(): Monday=0..Sunday=6
    return names[date.weekday()]

def _month_name(date: dt.datetime, lang: str) -> str:
    names = _MONTH_NAMES.get(lang) or _MONTH_NAMES.get("en")
    # month 1..12
    return names[date.month - 1]

_lang = configuration.conf.email_template.language if configuration.conf.email_template.language in ["en","fr","he"] else "en"
now = dt.datetime.now()
start = now - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days)

placeholders = SafeFormatDict({
    "date": now.strftime("%Y-%m-%d"),
    "day_name": _day_name(now, _lang),
    "day_number": now.strftime("%d"),
    "month_name": _month_name(now, _lang),
    "month_number": now.strftime("%m"),
    "year": now.strftime("%Y"),
    "start_date": start.strftime("%Y-%m-%d"),
    "start_day_name": _day_name(start, _lang),
    "start_day_number": start.strftime("%d"),
    "start_month_name": _month_name(start, _lang),
    "start_month_number": start.strftime("%m"),
    "start_year": start.strftime("%Y")
})