from source.configuration import conf, logging
import requests
import datetime as dt

def get_root_items():
    headers = {
        "Authorization": f'MediaBrowser Token="{conf.jellyfin.api_token}"'
    }

    response = requests.get(f'{conf.jellyfin.url}/Items', headers=headers)
    if response.status_code != 200:
        logging.error(f"Error while getting the root items, status code: {response.status_code}.")
        raise Exception(f"Error while getting the root items, status code: {response.status_code}. Answer: {response.text}.")
    return response.json()["Items"]


def get_item_from_parent(parent_id, type, minimum_creation_date=None):
    if type not in ["movie", "tv"]:
        raise Exception(f"Error while retrieving a media from Jellyfin. Type must be 'movie' or 'tv'. Got {type}")
    headers = {
        "Authorization": f'MediaBrowser Token="{conf.jellyfin.api_token}"'
    }



    response = requests.get(f'{conf.jellyfin.url}/Items?ParentId={parent_id}&fields=DateCreated,ProviderIds&Recursive=true', headers=headers)
    if response.status_code != 200:
        logging.error(f"Error while getting the items from parent, status code: {response.status_code}.")
        raise Exception(f"Error while getting the items from parent, status code: {response.status_code}. Answer: {response.text}.")
    if not minimum_creation_date:
        return response.json()["Items"], response.json()["TotalRecordCount"]
    else:
        recent_items = []
        for item in response.json()["Items"]:
            if (item.get("Type") == "Episode" and item.get("LocationType") == "Virtual") or (item.get("Type") == "Movie" and item.get("LocationType") == "Virtual"):
                # see https://github.com/SeaweedbrainCY/jellyfin-newsletter/issues/28 for context
                logging.debug(f"Skipping item because it is a virtual item. Item : {item}")
                continue
            creation_date = dt.datetime.strptime(item["DateCreated"].split("T")[0], "%Y-%m-%d")
            if creation_date > minimum_creation_date:
                logging.debug(f"Item {item['Name']} is more recent than {minimum_creation_date} (added on {creation_date}). Adding it to the list.")
                logging.debug("Item details: " + str(item))
                recent_items.append(item)
        return recent_items, response.json()["TotalRecordCount"]



def get_item_from_parent_by_id(parent_id, item_id):
    headers = {
        "Authorization": f'MediaBrowser Token="{conf.jellyfin.api_token}"'
    }
    response = requests.get(f'{conf.jellyfin.url}/Items?ParentId={parent_id}&fields=DateCreated,ProviderIds&Recursive=true&ids={item_id}', headers=headers)
    if response.status_code != 200:
        logging.error(f"Error while getting the items from parent, status code: {response.status_code}.")
        raise Exception(f"Error while getting the items from parent, status code: {response.status_code}. Answer: {response.text}.")
    for item in response.json()["Items"]:
        if "Id" in item.keys():
            if item["Id"] == item_id:
                return item