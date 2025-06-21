import sys 
from source import configuration, JellyfinAPI, TmdbAPI, email_template, email_controller
import datetime as dt
from source.configuration import logging



def populate_series_item(series_items, item):
    if "DateCreated" not in item.keys():
        logging.warning(f"Item {item['SeriesName']} has no creation date.")
        item["DateCreated"] = None
    if "ProductionYear" not in item.keys():
        logging.warning(f"Item {item['SeriesName']} has no production year.")
        item["ProductionYear"] = 0

    tmdb_id = None
    if "ProviderIds" in item.keys():
        if "Tmdb" in item["ProviderIds"].keys():
            tmdb_id = item["ProviderIds"]["Tmdb"]

    if tmdb_id is not None: # id provided by Jellyfin
        tmdb_info = TmdbAPI.get_media_detail_from_id(id=tmdb_id, type="tv")
    else:
        logging.info(f"Item {item['SeriesName']} has no TMDB id, searching by title.")
        tmdb_info = TmdbAPI.get_media_detail_from_title(title=item["SeriesName"], type="tv", year=item["ProductionYear"])

    if tmdb_info is None:
        logging.warning(f"Item {item['Name']} has not been found on TMDB. Skipping.")
    else:
        if item['Name'] not in series_items.keys():
            series_items[item['Name']] = {
                "is_serie_item_initialized": True,
            }
        series_items[item['Name']]["created_on"] = item["DateCreated"]
        if "seasons" not in series_items[item['Name']].keys():
            series_items[item['Name']]["seasons"] = []
        if "overview" not in tmdb_info.keys():
            logging.warning(f"Item {item['SeriesName']} has no overview.")
            tmdb_info["Overview"] = "No overview available."
        series_items[item['Name']]["description"] = tmdb_info["overview"]
        series_items[item['Name']]["year"] = item["ProductionYear"]
        series_items[item['Name']]["poster"] = f"https://image.tmdb.org/t/p/w500{tmdb_info['poster_path']}" if tmdb_info["poster_path"] else "https://redthread.uoregon.edu/files/original/affd16fd5264cab9197da4cd1a996f820e601ee4.png"

def catch_undefined_series(series_items, watched_tv_folders_id):
    """
    This script only search for items created within the defined past days. 
    If a series has been created before that period, but episodes have been created during this period, the series will appear with missing information in the newsletter.
    To catch this case, we will parse the final series_item array and verify that no series is missing information.
    """
    for folder_id in watched_tv_folders_id:
        for serie_name in series_items.keys():
            if not series_items[serie_name]["is_serie_item_initialized"]:
                item = JellyfinAPI.get_item_from_parent_by_name(parent_id=folder_id, name=serie_name)
                if item is not None:
                    populate_series_item(series_items, item)
                    




if __name__ == "__main__":
    logging.info("""

Jellyfin Newsletter is starting ....
##############################################


""")

    folders = JellyfinAPI.get_root_items()
    watched_film_folders_id = []
    watched_tv_folders_id = []
    for item in folders:
        if item["Name"] in configuration.conf.jellyfin.watched_film_folders :
           watched_film_folders_id.append(item["Id"])
        elif item["Name"] in configuration.conf.jellyfin.watched_tv_folders :
            watched_tv_folders_id.append(item["Id"])
        else:
            logging.warning(f"Folder {item['Name']} is not watched. Skipping.")

    total_movie = 0
    total_tv = 0
    movie_items = {}
    series_items = {}


    for folder_id in watched_film_folders_id:
        items, total_count = JellyfinAPI.get_item_from_parent(parent_id=folder_id,type="movie", minimum_creation_date=dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days))
        total_movie += total_count
        for item in items:
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

                movie_items[item["Name"]] = {
                    "year":item["ProductionYear"],
                    "created_on":item["DateCreated"],
                    "description": tmdb_info["overview"],
                    "poster": f"https://image.tmdb.org/t/p/w500{tmdb_info['poster_path']}" if tmdb_info["poster_path"] else "https://redthread.uoregon.edu/files/original/affd16fd5264cab9197da4cd1a996f820e601ee4.png"
                }
            
    
    for folder_id in watched_tv_folders_id:
        items, total_count = JellyfinAPI.get_item_from_parent(parent_id=folder_id, type="tv", minimum_creation_date=dt.datetime.now() - dt.timedelta(days=configuration.conf.jellyfin.observed_period_days))
        total_tv += total_count
        for item in items:
            if item["Type"] == "Series":
                populate_series_item(series_items, item)
            elif item["Type"] == "Episode":
                if (item['SeriesName'] in series_items.keys() and item["SeasonName"] not in series_items[item["SeriesName"]]["seasons"] ) or (item['SeriesName'] not in series_items.keys()):
                    if item['SeriesName'] not in series_items.keys():
                        series_items[item['SeriesName']] = {
                            "seasons": [],
                            "is_serie_item_initialized": False,
                            "created_on": "undefined",
                            "description": "No description available.",
                            "year": "undefined",
                            "poster": "https://redthread.uoregon.edu/files/original/affd16fd5264cab9197da4cd1a996f820e601ee4.png"
                        }
                    series_items[item['SeriesName']]["seasons"].append(item["SeasonName"])

            
    catch_undefined_series(series_items=series_items, watched_tv_folders_id=watched_tv_folders_id)
    template = email_template.populate_email_template(movies=movie_items, series=series_items, total_tv=total_tv, total_movie=total_movie)


    email_controller.send_email(template)

    logging.info("""All done here. Exiting.
    
    
##############################################
Thanks for using Jellyfin Newsletter!
Developed by Seaweedbrain, under MIT License.""")

    


