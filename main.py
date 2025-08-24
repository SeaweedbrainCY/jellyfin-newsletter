import sys 
from source import configuration, JellyfinAPI, TmdbAPI, email_template, email_controller
import datetime as dt
from source.configuration import logging
from source.configuration_checker import check_configuration
from apscheduler.schedulers.blocking import BlockingScheduler
from apscheduler.triggers.cron import CronTrigger
import source.utils as utils


def populate_series_item_from_episode(series_items, item):
    """
    Populate the series item with required information to build the email content. 
    It takes an episode, populate the serie item with the episode information, and add the episode to the series item.
    series_items format : 
    {
        "id": {
            "series_name": "SeriesName", # Name of the series, provided by Jellyfin
            "created_on": "2023-10-01T12:00:00Z", # Creation date of the series item, i.e. of added season or the added episode, or the series itfself
            "description": "This is a series description.", # Since episode rarely includes TMBD id, even if the episode is alone, we will use the series description.
            "year": 2023, # Production year of the series, provided by Jellyfin
            "poster": "", # Poster of the series, provided by TMDB
            "seasons": [""] # List of seasons in the series, provided by Jellyfin
            "episodes" : [""]
        }
    }
    """


    required_keys = ["SeriesId", "SeriesName", "SeasonName"]
    for key in required_keys:
        if key not in item.keys():
            logging.warning(f"Item {item} has no {key}. Skipping.")
            return
    if item["SeriesId"] not in series_items.keys():
        series_items[item["SeriesId"]] = {
            "series_name": item["SeriesName"],  # Name of the series, provided by Jellyfin
            "episodes": [],
            "seasons": [],
            "created_on": "undefined",
            "description": "No description available.",  # will be populated later, when parsing the series item
            "year": "undefined",# will be populated later, when parsing the series item
            "poster": "https://redthread.uoregon.edu/files/original/affd16fd5264cab9197da4cd1a996f820e601ee4.png"# will be populated later, when parsing the series item
        }
    if item["SeasonName"] not in series_items[item["SeriesId"]]["seasons"]:
        series_items[item["SeriesId"]]["seasons"].append(item["SeasonName"])
    series_items[item["SeriesId"]]["episodes"].append(item.get('IndexNumber'))
    if series_items[item["SeriesId"]]["created_on"] != "undefined" or series_items[item["SeriesId"]]["created_on"] is not None:
        try: 
            if dt.datetime.fromisoformat(series_items[item["SeriesId"]]["created_on"]) < dt.datetime.fromisoformat(item["DateCreated"]):
                series_items[item["SeriesId"]]["created_on"] = item["DateCreated"]
        except:
            pass
    series_items[item["SeriesId"]]["created_on"] = item.get("DateCreated", "undefined") 


def populate_series_item_with_series_related_information(series_items, watched_tv_folders_id):
    """
    populate_series_item_from_episode will populate the series item with the episode information, but it will not include the series information (description, year, poster).
    This function will populate the series item with the series information.
    """
    for folder_id in watched_tv_folders_id:
        for series_id in series_items.keys():
            item = JellyfinAPI.get_item_from_parent_by_id(parent_id=folder_id, item_id=series_id)
            if item is not None:
                if "Type" not in item.keys() or item["Type"] != "Series":
                    logging.warning(f"Item {item} is not a series. Skipping.")
                    continue
                required_keys = ["Name", "Id"]
                for key in required_keys:
                    if key not in item.keys():
                        logging.warning(f"Item {item} has no {key}. Skipping.")
                        continue
                series_items[item['Id']]["year"] = item["ProductionYear"]
                tmdb_id = None
                if "ProviderIds" in item.keys():
                    if "Tmdb" in item["ProviderIds"].keys():
                        tmdb_id = item["ProviderIds"]["Tmdb"]
    
                if tmdb_id is not None: # id provided by Jellyfin
                    try:
                        tmdb_info = TmdbAPI.get_media_detail_from_id(id=tmdb_id, type="tv")
                    except Exception as e:
                        logging.error(f"Item {item['Name']} could not be retrieved from TMDB by id due to an API error: {e}")     
                        logging.info(f"Retrying search for item {item} by title.")

                if tmdb_id is None or tmdb_info is None:
                    logging.info(f"Item {item} has no TMDB id or search by id failed. Searching by title.")
                    try:
                        tmdb_info = TmdbAPI.get_media_detail_from_title(title=item["Name"], type="tv", year=item["ProductionYear"])
                    except Exception as e:
                        logging.error(f"Item {item['Name']} could not be retrieved from TMDB by title due to an API error: {e}")   
                                       
                if tmdb_info is None:
                    logging.warning(f"Item {item['Name']} has not been found on TMDB. Skipping.")
                else:
                    if "overview" not in tmdb_info.keys():
                        logging.warning(f"Item {item['Name']} has no overview.")
                        tmdb_info["Overview"] = "No overview available."
                    series_items[item['Id']]["description"] = tmdb_info["overview"]
                    
                    series_items[item['Id']]["poster"] = f"https://image.tmdb.org/t/p/w500{tmdb_info['poster_path']}" if tmdb_info["poster_path"] else "https://redthread.uoregon.edu/files/original/affd16fd5264cab9197da4cd1a996f820e601ee4.png"
            else:
                logging.warning(f"Item {series_id} has not been found in Jellyfin. Skipping.")

    


def send_newsletter():
    logging.info("Sending newsletter ...")
    folders = JellyfinAPI.get_root_items()
    watched_film_folders_id = []
    watched_tv_folders_id = []
    for item in folders:
        if "Name" not in item:
            logging.warning(f"Item {item} has no Name. Skipping.")
            continue
        if item["Name"] in configuration.conf.jellyfin.watched_film_folders :
           watched_film_folders_id.append(item["Id"])
           logging.info(f"Folder {item['Name']} is watched for films.")
        elif item["Name"] in configuration.conf.jellyfin.watched_tv_folders :
            watched_tv_folders_id.append(item["Id"])
            logging.info(f"Folder {item['Name']} is watched for TV series.")
        else:
            logging.warning(f"Folder {item['Name']} is not watched. Skipping. Add \"{item['Name']}\" in your watched folder to include it.")

    total_movie = 0
    total_tv = 0
    movie_items = {}
    series_items = {}


    for folder_id in watched_film_folders_id:
        items, total_count = JellyfinAPI.get_item_from_parent(parent_id=folder_id,type="movie", minimum_creation_date=dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days))
        total_movie += total_count
        for item in items:
            required_keys = ["Name", "Id", "DateCreated"]
            for key in required_keys:
                if key not in item.keys():
                    logging.warning(f"Item {item} has no {key}. Skipping.")
                    continue
            if configuration.conf.jellyfin.ignore_item_added_before_last_newsletter:
                last_newsletter_date = utils.get_last_newsletter_date()
                if last_newsletter_date is not None:
                    if item["DateCreated"] is not None and dt.datetime.strptime(item["DateCreated"].split("T")[0], "%Y-%m-%d") < last_newsletter_date:
                        logging.info(f"ignore_item_added_before_last_newsletter is set to True and Item {item['Name']} was added before the last newsletter. Ignoring.")
                        continue
            tmdb_id = None
            if "ProductionYear"  not in item.keys():
                logging.warning(f"Item {item['Name']} has no production year.")
                item["ProductionYear"] = 0
            if "DateCreated" not in item.keys():
                logging.warning(f"Item {item['Name']} has no creation date.")
                item["DateCreated"] = None
            if "ProviderIds" in item.keys():
                if "Tmdb" in item["ProviderIds"].keys():
                    tmdb_id = item["ProviderIds"]["Tmdb"]
            
            
            if tmdb_id is not None: # id provided by Jellyfin
                tmdb_info = TmdbAPI.get_media_detail_from_id(id=tmdb_id, type="movie")
            else:
                logging.info(f"Item {item['Name']} has no TMDB id, searching by title.")
                tmdb_info = TmdbAPI.get_media_detail_from_title(title=item["Name"], type="movie", year=item["ProductionYear"])

            if tmdb_info is None:
                logging.warning(f"Item {item['Name']} has not been found on TMDB. Skipping.")
            else:
                if "overview" not in tmdb_info.keys():
                    logging.warning(f"Item {item['Name']} has no overview.")
                    tmdb_info["overview"] = "No overview available."

                movie_items[item["Id"]] = {
                    "name": item["Name"],
                    "year":item["ProductionYear"],
                    "created_on":item["DateCreated"],
                    "description": tmdb_info["overview"],
                    "poster": f"https://image.tmdb.org/t/p/w500{tmdb_info['poster_path']}" if tmdb_info["poster_path"] else "https://redthread.uoregon.edu/files/original/affd16fd5264cab9197da4cd1a996f820e601ee4.png"
                }
            
    
    for folder_id in watched_tv_folders_id:
        items, total_count = JellyfinAPI.get_item_from_parent(parent_id=folder_id, type="tv", minimum_creation_date=dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days))
        total_tv += total_count
        for item in items:
            if configuration.conf.jellyfin.ignore_item_added_before_last_newsletter:
                last_newsletter_date = utils.get_last_newsletter_date()
                if last_newsletter_date is not None:
                    if item["DateCreated"] is not None and dt.datetime.strptime(item["DateCreated"].split("T")[0], "%Y-%m-%d") < last_newsletter_date:
                        logging.info(f"ignore_item_added_before_last_newsletter is set to True and Item {item.get('Name')} was added before the last newsletter. Ignoring.")
                        continue
            if item["Type"] == "Episode":
                populate_series_item_from_episode(series_items, item)
    
            
    populate_series_item_with_series_related_information(series_items=series_items, watched_tv_folders_id=watched_tv_folders_id)
    logging.debug("Series populated : " + str(series_items))
    if len(movie_items) + len(series_items) > 0:
        template = email_template.populate_email_template(movies=movie_items, series=series_items, total_tv=total_tv, total_movie=total_movie)

        # Use the new send_newsletter function with preview/dry-run support
        result = email_controller.send_newsletter(
            html_content=template,
            movies=movie_items,
            series=series_items,
            total_tv=total_tv,
            total_movie=total_movie
        )

        # Log results based on mode
        if result["mode"] == "normal":
            logging.info(f"All emails sent successfully to {result['sent_count']} recipients.")
        elif result["mode"] == "preview":
            logging.info(f"Preview generated successfully (preview-only mode).")
        elif result["mode"] == "dry-run":
            smtp_status = "PASSED" if result["smtp_tested"] else "FAILED"
            logging.info(f"Dry-run completed. SMTP test: {smtp_status}")
        
        logging.info("Newsletter processing completed.")
    else:
        logging.warning("No new items found in watched folders. No email sent.")
    
    logging.info("""


##############################################
Newsletter sent. 
Thanks for using Jellyfin Newsletter!
Developed by Seaweedbrain, under MIT License.""")





def newsletter_job():
    """
    Used to run the newsletter, called by scheduler. 
    Used to handle exceptions and logging.
    """
    try:
        send_newsletter()
    except Exception as e:
        logging.error(f"[FATAL] An error occurred while sending the newsletter: {e}")
        logging.error("Sending newsletter failed. Program will continue to run and retry at the next scheduled time.")



if __name__ == "__main__":
    current_version = ""
    try:
        with open("VERSION", "r") as version_file:
            current_version = version_file.read().strip()
    except :
       current_version = "unknown version"        
        
    logging.info(f"""

Jellyfin Newsletter {current_version} is starting ....
##############################################



""")
    logging.info("Checking configuration ...")
    try:
        check_configuration()
    except Exception as e:
        logging.error(f"[FATAL] Configuration check failed: {e}")
        sys.exit(1)
    logging.info("Configuration check passed.")

    if configuration.conf.scheduler.enabled:
        try:
            scheduler = BlockingScheduler()
            trigger = CronTrigger().from_crontab(configuration.conf.scheduler.cron)
        except Exception as e:
            logging.error(f"[FATAL] Failed to initialize scheduler: {e}")
            sys.exit(1)

        scheduler.add_job(newsletter_job, trigger)
        logging.info(f"Newsletter scheduler started. Next run at {trigger.get_next_fire_time(None, dt.datetime.now()).isoformat()}")
        scheduler.start()
        
    else:
        logging.info("Scheduler is disabled. Newsletter will run once, now.")
        send_newsletter()

        







    

    


