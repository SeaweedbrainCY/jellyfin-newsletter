from source import configuration, context, utils
from source.language_utils import TRANSLATIONS
import re

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
    with open(f"./themes/new_media/{configuration.conf.email_template.theme}/main.html", encoding='utf-8') as template_file:
        template = template_file.read()
        
        if configuration.conf.email_template.language in configuration.conf.email_template.available_lang:
            for key in TRANSLATIONS[configuration.conf.email_template.language]:
                template = re.sub(
                    r"\${" + key + "}", 
                    TRANSLATIONS[configuration.conf.email_template.language][key], 
                    template
                )
            missing = re.findall(r"\${([^}]+)}", template)
            if missing:
                configuration.logging.info(f"Missing translations for language {configuration.conf.email_template.language}: {', '.join(set(missing))} ")
        else:
            raise Exception(f"[FATAL] Language {configuration.conf.email_template.language} not supported. Supported languages are {', '.join(configuration.conf.email_template.available_lang)}.")

        # lang/dir for the root HTML tag
        html_lang = configuration.conf.email_template.language if configuration.conf.email_template.language in configuration.conf.email_template.available_lang else "en"
        text_dir = "rtl" if configuration.conf.email_template.language == "he" else "ltr"

        # Wrap English link texts in <bdi> for Hebrew to preserve word order
        project_text = "Jellyfin Newsletter"
        if configuration.conf.email_template.language == "he":
            project_text = f"<bdi>{project_text}</bdi>"

        custom_keys = [
            {"key": "title", "value": configuration.conf.email_template.title.format_map(context.placeholders)}, 
            {"key": "subtitle", "value": configuration.conf.email_template.subtitle.format_map(context.placeholders)},
            {"key": "jellyfin_url", "value": configuration.conf.email_template.jellyfin_url},
            {"key": "jellyfin_owner_name", "value": configuration.conf.email_template.jellyfin_owner_name.format_map(context.placeholders)},
            {"key": "unsubscribe_email", "value": configuration.conf.email_template.unsubscribe_email.format_map(context.placeholders)},
            {"key": "html_lang", "value": html_lang},
            {"key": "dir", "value": text_dir},
            {"key": "project_link_text", "value": project_text},
        ]
        
        for key in custom_keys:
            template = re.sub(r"\${" + key["key"] + "}", key["value"], template)

        # Movies section
        if movies:
            template = re.sub(r"\${display_movies}", "", template)
            movies_html = ""
            sort_mode = getattr(configuration.conf.email_template, "sort_mode", "date_asc")
            if sort_mode == "name_asc":
                movie_items_sorted = sorted(movies.items(), key=lambda kv: kv[1]["name"].casefold())
            elif sort_mode == "name_desc":
                movie_items_sorted = sorted(movies.items(), key=lambda kv: kv[1]["name"].casefold(), reverse=True)
            elif sort_mode == "date_desc":
                movie_items_sorted = sorted(movies.items(), key=lambda kv: kv[1]["created_on"] or "", reverse=True)
            else:  # date_asc
                movie_items_sorted = sorted(movies.items(), key=lambda kv: kv[1]["created_on"] or "")

            for movie_id, movie_data in movie_items_sorted:
                added_date = movie_data["created_on"].split("T")[0]
                added_date_html = f"<bdi>{added_date}</bdi>"
                movie_overview_style = "display: none;"
                if include_overview:
                    movie_overview_style = "display: block;"
                with open(f"./themes/new_media/{configuration.conf.email_template.theme}/movie.html", encoding='utf-8') as movie_template_file:
                    movie_template = movie_template_file.read()
                    movie_template = re.sub(r"\${movie_overview_style}", movie_overview_style, movie_template)
                    movie_template = re.sub(r"\${movie_poster}", movie_data['poster'], movie_template)
                    movie_template = re.sub(r"\${movie_name}", movie_data['name'], movie_template)
                    movie_template = re.sub(r"\${movie_added_on_label}", TRANSLATIONS[configuration.conf.email_template.language]['added_on'], movie_template)
                    movie_template = re.sub(r"\${movie_added_on}", added_date_html, movie_template)
                    movie_template = re.sub(r"\${movie_overview}", movie_data['description'], movie_template)
                    movies_html += movie_template

                
            template = re.sub(r"\${films}", movies_html, template)
        else:
            template = re.sub(r"\${display_movies}", "display:none", template)

        # TV Shows section
        if series:
            template = re.sub(r"\${display_tv}", "", template)
            series_html = ""
            sort_mode = getattr(configuration.conf.email_template, "sort_mode", "date_asc")
            if sort_mode == "name_asc":
                series_items_sorted = sorted(series.items(), key=lambda kv: kv[1]["series_name"].casefold())
            elif sort_mode == "name_desc":
                series_items_sorted = sorted(series.items(), key=lambda kv: kv[1]["series_name"].casefold(), reverse=True)
            elif sort_mode == "date_desc":
                series_items_sorted = sorted(series.items(), key=lambda kv: kv[1]["created_on"] or "", reverse=True)
            else:  # date_asc
                series_items_sorted = sorted(series.items(), key=lambda kv: kv[1]["created_on"] or "")

            for serie_id, serie_data in series_items_sorted:
                added_date = serie_data["created_on"].split("T")[0]
                added_date_html = f"<bdi>{added_date}</bdi>"
                if len(serie_data["seasons"]) == 1 :
                    if len(serie_data["episodes"]) == 1:
                        added_items_str = f"{serie_data['seasons'][0]}, {TRANSLATIONS[configuration.conf.email_template.language]['episode']} {serie_data['episodes'][0]}"
                    else:
                        episodes_ranges = utils.summarize_ranges(serie_data["episodes"])
                        if episodes_ranges is None:
                            added_items_str = f"{serie_data['seasons'][0]}, {TRANSLATIONS[configuration.conf.email_template.language]['new_episodes']}."
                        if len(episodes_ranges) == 1:
                            added_items_str = f"{serie_data['seasons'][0]}, {TRANSLATIONS[configuration.conf.email_template.language]['episodes']} {episodes_ranges[0]}"
                        else:
                            added_items_str = f"{serie_data['seasons'][0]}, {TRANSLATIONS[configuration.conf.email_template.language]['episodes']} {', '.join(episodes_ranges[:-1])} & {episodes_ranges[-1]}"
                else:
                    serie_data["seasons"].sort()
                    added_items_str = ", ".join(serie_data["seasons"])

                tv_title = f"<bdi>{serie_data['series_name']}: {added_items_str}</bdi>"
                tv_overview_style = "display: none;"
                if include_overview:
                    tv_overview_style = "display: block;"
                added_items_html = f"<bdi>{added_items_str}</bdi>"
                with open(f"./themes/new_media/{configuration.conf.email_template.theme}/tv.html", encoding='utf-8') as movie_template_file:
                    tv_template = movie_template_file.read()
                    tv_template = re.sub(r"\${tv_title}", serie_data['series_name'], tv_template)
                    tv_template = re.sub(r"\${tv_overview_style}", tv_overview_style, tv_template)
                    tv_template = re.sub(r"\${tv_overview}", serie_data['description'], tv_template)
                    tv_template = re.sub(r"\${tv_added_on}", added_date_html, tv_template)
                    tv_template = re.sub(r"\${tv_added_on_label}", TRANSLATIONS[configuration.conf.email_template.language]['added_on'], tv_template)
                    tv_template = re.sub(r"\${tv_poster}", serie_data['poster'], tv_template)
                    series_html += tv_template
                
                
            template = re.sub(r"\${tvs}", series_html, template)
        else:
            template = re.sub(r"\${display_tv}", "display:none", template)

        # Always use BDI for numbers to ensure proper display in mixed-direction content
        series_count_value = f"<bdi>{total_tv}</bdi>"
        movies_count_value = f"<bdi>{total_movie}</bdi>"
        template = re.sub(r"\${series_count}", series_count_value, template)
        template = re.sub(r"\${movies_count}", movies_count_value, template)

        return template