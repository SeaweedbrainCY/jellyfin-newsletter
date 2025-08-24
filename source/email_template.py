from source import configuration, context, utils
import re

translation = {
    "en":{
        "discover_now": "Discover now",
        "new_film": "New movies:",
        "new_tvs": "New shows:",
        "currently_available": "Currently available in Jellyfin:",
        "movies_label": "Movies",
        "episodes_label": "Episodes",
        "footer_label":"You are recieving this email because you are using ${jellyfin_owner_name}'s Jellyfin server. If you want to stop receiving these emails, you can unsubscribe by notifying ${unsubscribe_email}.",
        "added_on": "Added on",
        "episodes": "Episodes",
         "episode": "Episode",
         "new_episodes": "new episodes",
    },
    "fr":{
        "discover_now": "Découvrir maintenant",
        "new_film": "Nouveaux films :",
        "new_tvs": "Nouvelles séries :",
        "currently_available": "Actuellement disponible sur Jellyfin :",
        "movies_label": "Films",
        "episodes_label": "Épisodes",
        "footer_label":"Vous recevez cet email car vous utilisez le serveur Jellyfin de ${jellyfin_owner_name}. Si vous ne souhaitez plus recevoir ces emails, vous pouvez vous désinscrire en notifiant ${unsubscribe_email}.",
        "added_on": "Ajouté le",
        "episodes": "Épisodes",
        "episode": "Épisode",
        "new_episodes": "nouveaux épisodes",
    }
}

def populate_email_template(movies, series, total_tv, total_movie) -> str:
    include_overview = True
    if configuration.conf.email_template.display_overview_max_items == -1:
        include_overview = False
        configuration.logging.debug("display_overview_max_items is -1, overviews will not be included in the email template.")
    elif configuration.conf.email_template.display_overview_max_items == 0:
        include_overview = True
        configuration.logging.debug("display_overview_max_items is 0, overviews will  be included in the email template, no matter their number.")
    elif len(movies) + len(series) > configuration.conf.email_template.display_overview_max_items :
        include_overview = False
        configuration.logging.info(f"There are more than {configuration.conf.email_template.display_overview_max_items} new items, overview will not be included in the email template to avoid too much content.")
    with open("./template/new_media_notification.html", encoding='utf-8') as template_file:
        template = template_file.read()
        
        if configuration.conf.email_template.language in ["fr", "en"]:
            for key in translation[configuration.conf.email_template.language]:
                template = re.sub(
                    r"\${" + key + "}", 
                    translation[configuration.conf.email_template.language][key], 
                    template
                )
        else:
            raise Exception(f"[FATAL] Language {configuration.conf.email_template.language} not supported. Supported languages are fr and en")

        custom_keys = [
            {"key": "title", "value": configuration.conf.email_template.title.format_map(context.placeholders)}, 
            {"key": "subtitle", "value": configuration.conf.email_template.subtitle.format_map(context.placeholders)},
            {"key": "jellyfin_url", "value": configuration.conf.email_template.jellyfin_url},
            {"key": "jellyfin_owner_name", "value": configuration.conf.email_template.jellyfin_owner_name.format_map(context.placeholders)},
            {"key": "unsubscribe_email", "value": configuration.conf.email_template.unsubscribe_email.format_map(context.placeholders)}
        ]
        
        for key in custom_keys:
            template = re.sub(r"\${" + key["key"] + "}", key["value"], template)

        # Movies section
        if movies:
            template = re.sub(r"\${display_movies}", "", template)
            movies_html = ""
            
            for movie_id, movie_data in movies.items():
                added_date = movie_data["created_on"].split("T")[0]
                item_overview_html = ""
                if include_overview:
                    item_overview_html = f"""
<div class="movie-description" style="color: #dddddd !important; font-size: 14px !important; line-height: 1.4 !important;">
                                            {movie_data['description']}
</div>
"""
                movies_html += f"""
                <div class="movie_container" style="margin-bottom: 15px;">
                    <div class="movie_bg" style="background: url('{movie_data['poster']}') no-repeat center center; background-size: cover; border-radius: 10px;">
                        <table class="movie" width="100%" role="presentation" cellpadding="0" cellspacing="0" style="background: rgba(0, 0, 0, 0.7); border-radius: 10px; width: 100%;">
                            <tr>
                                <td class="movie-image" valign="middle" style="padding: 15px; text-align: center; width: 120px;">
                                    <img src="{movie_data['poster']}" alt="{movie_data['name']}" style="max-width: 100px; height: auto; display: block; margin: 0 auto;">
                                </td>
                                <td class="movie-content-cell" valign="middle" style="padding: 15px;">
                                    <div class="mobile-text-container">
                                        <h3 class="movie-title" style="color: #ffffff !important; margin: 0 0 5px !important; font-size: 18px !important;">{movie_data['name']}</h3>
                                        <div class="movie-date" style="color: #dddddd !important; font-size: 14px !important; margin: 0 0 10px !important;">
                                            {translation[configuration.conf.email_template.language]['added_on']} {added_date}
                                        </div>
                                        {item_overview_html}
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </div>
                </div>
                """
                
            template = re.sub(r"\${films}", movies_html, template)
        else:
            template = re.sub(r"\${display_movies}", "display:none", template)

        # TV Shows section
        if series:
            template = re.sub(r"\${display_tv}", "", template)
            series_html = ""
            
            for serie_id, serie_data in series.items():
                added_date = serie_data["created_on"].split("T")[0]
                if len(serie_data["seasons"]) == 1 :
                    if len(serie_data["episodes"]) == 1:
                        added_items_str = f"{serie_data['seasons'][0]}, {translation[configuration.conf.email_template.language]['episode']} {serie_data['episodes'][0]}"
                    else:
                        episodes_ranges = utils.summarize_ranges(serie_data["episodes"])
                        if episodes_ranges is None:
                            added_items_str = f"{serie_data['seasons'][0]}, {translation[configuration.conf.email_template.language]['new_episodes']}."
                        if len(episodes_ranges) == 1:
                            added_items_str = f"{serie_data['seasons'][0]}, {translation[configuration.conf.email_template.language]['episodes']} {episodes_ranges[0]}"
                        else:
                            added_items_str = f"{serie_data['seasons'][0]}, {translation[configuration.conf.email_template.language]['episodes']} {', '.join(episodes_ranges[:-1])} & {episodes_ranges[-1]}"
                else:
                    serie_data["seasons"].sort()
                    added_items_str = ", ".join(serie_data["seasons"])

                item_overview_html = ""
                if include_overview:
                    item_overview_html = f"""
<div class="movie-description" style="color: #dddddd !important; font-size: 14px !important; line-height: 1.4 !important;">
                                            {serie_data['description']}
                                        </div>
"""
                series_html += f"""
                <div class="movie_container" style="margin-bottom: 15px;">
                    <div class="movie_bg" style="background: url('{serie_data['poster']}') no-repeat center center; background-size: cover; border-radius: 10px;">
                        <table class="movie" width="100%" role="presentation" cellpadding="0" cellspacing="0" style="background: rgba(0, 0, 0, 0.7); border-radius: 10px; width: 100%;">
                            <tr>
                                <td class="movie-image" valign="middle" style="padding: 15px; text-align: center; width: 120px;">
                                    <img src="{serie_data['poster']}" alt="{serie_data['series_name']}" style="max-width: 100px; height: auto; display: block; margin: 0 auto;">
                                </td>
                                <td class="movie-content-cell" valign="middle" style="padding: 15px;">
                                    <div class="mobile-text-container">
                                        <h3 class="movie-title" style="color: #ffffff !important; margin: 0 0 5px !important; font-size: 18px !important;">{serie_data['series_name']}: {added_items_str}</h3>
                                        <div class="movie-date" style="color: #dddddd !important; font-size: 14px !important; margin: 0 0 10px !important;">
                                            {translation[configuration.conf.email_template.language]['added_on']} {added_date}
                                        </div>
                                        {item_overview_html}
                                    </div>
                                </td>
                            </tr>
                        </table>
                    </div>
                </div>
                """
                
            template = re.sub(r"\${tvs}", series_html, template)
        else:
            template = re.sub(r"\${display_tv}", "display:none", template)

        # Statistics section
        template = re.sub(r"\${series_count}", str(total_tv), template)
        template = re.sub(r"\${movies_count}", str(total_movie), template)
        
        return template