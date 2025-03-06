import requests
from source import configuration 
import json
import logging




def get_media_detail_from_title(title, type, year=None):
    year_query = f"&year={year}" if year else ""
    if type not in ["movie", "tv"]:
        logging.error(f"Error while retrieving a media from TMDB. Type must be 'movie' or 'tv'. Got {type}")
        return None
    url = f"https://api.themoviedb.org/3/search/{type}?query={title}&language=fr-fr{year_query}"

    headers = {
        "accept": "application/json",
        "Authorization": f"Bearer {configuration.conf.tmdb.api_key}"
    }

    response = requests.get(url, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Error while getting the token, status code: {response.status_code}. Answer: {response.json()}.")
    if response.json()["total_results"] == 1:
        return response.json()["results"][0]
    elif response.json()["total_results"] > 1:
        logging.warning(f"Warning, multiple results found for the title {title}.")
        max_popularity = 0
        best_result = None
        for result in response.json()["results"]:
            if result["popularity"] > max_popularity:
                max_popularity = result["popularity"]
                best_result = result
        logging.info(f"Returning the best result. Movie id {best_result["id"]} has been selected based on popularity.")
        return best_result
    else:
        logging.error(f"Error, no result found for the title {title}.")
        return None