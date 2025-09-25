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
         "footer_project_open_source": "is an open source project.",
         "footer_developed_by": "Developed with ❤️ by",
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
    "footer_project_open_source": "est un projet open source.",
    "footer_developed_by": "Développé avec ❤️ par",
    },
    "he":{
        "discover_now": "גלה עכשיו",
        "new_film": "סרטים חדשים:\u200f",
        "new_tvs": "סדרות חדשות:\u200f",
        # Add RLM (\u200f) after colon to keep it properly positioned in RTL text
        "currently_available": "זמין כעת בג'ליפין:\u200f",
        "movies_label": "סרטים",
        "episodes_label": "פרקים",
        "footer_label":"אתם מקבלים מייל זה משום שאתם משתמשים בשרת ג'ליפין של ${jellyfin_owner_name}. כדי להפסיק לקבל מיילים אלה, ניתן לבקש להסיר ב־${unsubscribe_email}.",
        "added_on": "נוסף בתאריך",
        "episodes": "פרקים",
        "episode": "פרק",
        "new_episodes": "פרקים חדשים",
        "footer_project_open_source": "הוא פרויקט קוד פתוח.",
        "footer_developed_by": "פותח באהבה על ידי",
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
        
        if configuration.conf.email_template.language in ["fr", "en", "he"]:
            for key in translation[configuration.conf.email_template.language]:
                template = re.sub(
                    r"\${" + key + "}", 
                    translation[configuration.conf.email_template.language][key], 
                    template
                )
        else:
            raise Exception(f"[FATAL] Language {configuration.conf.email_template.language} not supported. Supported languages are fr, en and he")

        # lang/dir for the root HTML tag
        html_lang = configuration.conf.email_template.language if configuration.conf.email_template.language in ["en","fr","he"] else "en"
        text_dir = "rtl" if configuration.conf.email_template.language == "he" else "ltr"

        # Wrap English link texts in <bdi> for Hebrew to preserve word order
        project_text = "Jellyfin Newsletter"
        developer_text = "Seaweedbrain"
        if configuration.conf.email_template.language == "he":
            project_text = f"<bdi>{project_text}</bdi>"
            developer_text = f"<bdi>{developer_text}</bdi>"

        custom_keys = [
            {"key": "title", "value": configuration.conf.email_template.title.format_map(context.placeholders)}, 
            {"key": "subtitle", "value": configuration.conf.email_template.subtitle.format_map(context.placeholders)},
            {"key": "jellyfin_url", "value": configuration.conf.email_template.jellyfin_url},
            {"key": "jellyfin_owner_name", "value": configuration.conf.email_template.jellyfin_owner_name.format_map(context.placeholders)},
            {"key": "unsubscribe_email", "value": configuration.conf.email_template.unsubscribe_email.format_map(context.placeholders)},
            {"key": "html_lang", "value": html_lang},
            {"key": "dir", "value": text_dir},
            {"key": "project_link_text", "value": project_text},
            {"key": "developer_link_text", "value": developer_text},
        ]
        
        for key in custom_keys:
            template = re.sub(r"\${" + key["key"] + "}", key["value"], template)

        # Movies section
        if movies:
            template = re.sub(r"\${display_movies}", "", template)
            movies_html = ""
            
            for movie_id, movie_data in movies.items():
                added_date = movie_data["created_on"].split("T")[0]
                # Always use BDI for dates to ensure proper display in mixed-direction content
                added_date_html = f"<bdi>{added_date}</bdi>"
                item_overview_html = ""
                if include_overview:
                    item_overview_html = f"""
<div class="movie-description" style="color: #dddddd !important; font-size: 14px !important; line-height: 1.4 !important;">
                                            {movie_data['description']}
</div>
"""
                image_cell = f"""
                                <td class=\"movie-image\" valign=\"middle\" style=\"padding: 15px; text-align: center; width: 120px;\"> 
                                    <img src=\"{movie_data['poster']}\" alt=\"{movie_data['name']}\" style=\"max-width: 100px; height: auto; display: block; margin: 0 auto;\">
                                </td>
                """
                content_cell = f"""
                                <td class=\"movie-content-cell\" valign=\"middle\" style=\"padding: 15px;\">
                                    <div class=\"mobile-text-container\">
                                        <h3 class=\"movie-title\" style=\"color: #ffffff !important; margin: 0 0 5px !important; font-size: 18px !important;\">{movie_data['name']}</h3>
                                        <div class=\"movie-date\" style=\"color: #dddddd !important; font-size: 14px !important; margin: 0 0 10px !important;\">
                                            {translation[configuration.conf.email_template.language]['added_on']} {added_date_html}
                                        </div>
                                        {item_overview_html}
                                    </div>
                                </td>
                """
                movies_html += f"""
                <div class=\"movie_container\" style=\"margin-bottom: 15px;\">
                    <div class=\"movie_bg\" style=\"background: url('{movie_data['poster']}') no-repeat center center; background-size: cover; border-radius: 10px;\">
                        <table class=\"movie\" width=\"100%\" role=\"presentation\" cellpadding=\"0\" cellspacing=\"0\" style=\"background: rgba(0, 0, 0, 0.7); border-radius: 10px; width: 100%;\">
                            <tr class=\"movie-row\">
{image_cell}
{content_cell}
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
                # Always use BDI for dates to ensure proper display in mixed-direction content
                added_date_html = f"<bdi>{added_date}</bdi>"
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

                # Always use BDI for mixed-direction content
                added_items_str_display = f"<bdi>{added_items_str}</bdi>"

                item_overview_html = ""
                if include_overview:
                    item_overview_html = f"""
<div class="movie-description" style="color: #dddddd !important; font-size: 14px !important; line-height: 1.4 !important;">
                                            {serie_data['description']}
                                        </div>
"""
                s_image_cell = f"""
                                <td class=\"movie-image\" valign=\"middle\" style=\"padding: 15px; text-align: center; width: 120px;\">
                                    <img src=\"{serie_data['poster']}\" alt=\"{serie_data['series_name']}\" style=\"max-width: 100px; height: auto; display: block; margin: 0 auto;\">
                                </td>
                """
                s_content_cell = f"""
                                <td class=\"movie-content-cell\" valign=\"middle\" style=\"padding: 15px;\">
                                    <div class=\"mobile-text-container\"> 
                                        <h3 class=\"movie-title\" style=\"color: #ffffff !important; margin: 0 0 5px !important; font-size: 18px !important;\">{serie_data['series_name']}: {added_items_str_display}</h3>
                                        <div class=\"movie-date\" style=\"color: #dddddd !important; font-size: 14px !important; margin: 0 0 10px !important;\">
                                            {translation[configuration.conf.email_template.language]['added_on']} {added_date_html}
                                        </div>
                                        {item_overview_html}
                                    </div>
                                </td>
                """
                series_html += f"""
                <div class=\"movie_container\" style=\"margin-bottom: 15px;\">
                    <div class=\"movie_bg\" style=\"background: url('{serie_data['poster']}') no-repeat center center; background-size: cover; border-radius: 10px;\">
                        <table class=\"movie\" width=\"100%\" role=\"presentation\" cellpadding=\"0\" cellspacing=\"0\" style=\"background: rgba(0, 0, 0, 0.7); border-radius: 10px; width: 100%;\">
                            <tr class=\"movie-row\">
{s_image_cell}
{s_content_cell}
                            </tr>
                        </table>
                    </div>
                </div>
                """
                
            template = re.sub(r"\${tvs}", series_html, template)
        else:
            template = re.sub(r"\${display_tv}", "display:none", template)

        # Always use BDI for numbers to ensure proper display in mixed-direction content
        series_count_value = f"<bdi>{total_tv}</bdi>"
        movies_count_value = f"<bdi>{total_movie}</bdi>"
        template = re.sub(r"\${series_count}", series_count_value, template)
        template = re.sub(r"\${movies_count}", movies_count_value, template)

        return template