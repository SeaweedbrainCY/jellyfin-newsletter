from sqlalchemy import Column, Integer, String, DateTime
from database import Base
import datetime


class MediaNotificationHistory(Base):
    __tablename__ = "media_notification_history"

    id = Column(Integer, primary_key=True, index=True)
    media_id = Column(String, index=True)
    media_name = Column(String, index=True)
    media_type = Column(String, index=True) 
    timestamp = Column(DateTime, default=datetime.datetime.utcnow)