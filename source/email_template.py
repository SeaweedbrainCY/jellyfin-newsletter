from source import configuration, context, utils
import re

# While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
TRANSLATIONS = {
    "ca": {
        "discover_now": "Descobreix ara",
        "new_film": "Pel·lícules noves:",
        "new_tvs": "Sèries noves:",
        "currently_available": "Disponible actualment a Jellyfin:",
        "movies_label": "Pel·lícules",
        "episodes_label": "Episodis",
        "footer_label": "Rebeu aquest correu electrònic perquè utilitzeu el servidor Jellyfin de ${jellyfin_owner_name}. Si no voleu rebre més aquests correus, podeu donar-vos de baixa notificant-ho a ${unsubscribe_email}.",
        "added_on": "Afegit el",
        "episodes": "Episodis",
        "episode": "Episodi",
        "new_episodes": "episodis nous",
        "footer_project_open_source": "és un projecte de codi obert.",        
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Desenvolupat amb ❤️ per <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> i <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">els col·laboradors</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, amb llicència AGPLv3.",
    },
    "de": {
        "discover_now": "Jetzt entdecken",
        "new_film": "Neue Filme:",
        "new_tvs": "Neue Serien:",
        "currently_available": "Derzeit verfügbar in Jellyfin:",
        "movies_label": "Filme",
        "episodes_label": "Episoden",
        "footer_label": "Sie erhalten diese E-Mail, da Sie den Jellyfin-Server von ${jellyfin_owner_name} nutzen. Wenn Sie diese E-Mails nicht mehr erhalten möchten, können Sie sich abmelden, indem Sie ${unsubscribe_email} benachrichtigen.",
        "added_on": "Hinzugefügt am",
        "episodes": "Episoden",
        "episode": "Episode",
        "new_episodes": "neue Episoden",
        "footer_project_open_source": "ist ein Open-Source-Projekt.",
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Entwickelt mit ❤️ von <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> und <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">Mitwirkenden</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, lizenziert unter AGPLv3.",
    },
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
         # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
         "footer_developed_by": 'Developed with ❤️ by <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> and  <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">contributors</a>.',
         "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, licensed under AGPLv3.",
    },
    "es": {
        "discover_now": "Descubrir ahora",
        "new_film": "Nuevas películas:",
        "new_tvs": "Nuevas series:",
        "currently_available": "Disponible actualmente en Jellyfin:",
        "movies_label": "Películas",
        "episodes_label": "Episodios",
        "footer_label": "Recibes este correo electrónico porque utilizas el servidor Jellyfin de ${jellyfin_owner_name}. Si no deseas seguir recibiendo estos correos, puedes darte de baja notificándolo a ${unsubscribe_email}.",
        "added_on": "Añadido el",
        "episodes": "Episodios",
        "episode": "Episodio",
        "new_episodes": "nuevos episodios",
        "footer_project_open_source": "es un proyecto de código abierto.",
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Desarrollado con ❤️ por <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> y <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">los colaboradores</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, bajo licencia AGPLv3.",
    },
    "fi": {
        "discover_now": "Tutustu nyt",
        "new_film": "Uudet elokuvat:",
        "new_tvs": "Uudet sarjat:",
        "currently_available": "Tällä hetkellä saatavilla Jellyfinissä:",
        "movies_label": "Elokuvat",
        "episodes_label": "Jaksot",
        "footer_label": "Saat tämän sähköpostin, koska käytät ${jellyfin_owner_name}:n Jellyfin-palvelinta. Jos haluat lopettaa näiden sähköpostien vastaanottamisen, voit peruuttaa tilauksen ilmoittamalla osoitteeseen ${unsubscribe_email}.",
        "added_on": "Lisätty",
        "episodes": "Jaksot",
        "episode": "Jakso",
        "new_episodes": "uudet jaksot",
        "footer_project_open_source": "on avoimen lähdekoodin projekti.",
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Kehitetty ❤️:lla <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> ja <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">yhteistyökumppanit</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, lisensoitu AGPLv3-lisenssillä.",
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
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Développé avec ❤️ par <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> et <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">les contributeurs</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, sous licence AGPLv3.",
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
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'פותח באהבה ❤️ על ידי <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> ו<a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">תורמים</a>.',
        "license_and_copyright": "זכויות יוצרים © 2025 Nathan Stchepinsky, ברישיון AGPLv3.", 
    },
    "it": {
        "discover_now": "Scopri ora",
        "new_film": "Nuovi film:",
        "new_tvs": "Nuove serie:",
        "currently_available": "Attualmente disponibile su Jellyfin:",
        "movies_label": "Film",
        "episodes_label": "Episodi",
        "footer_label": "Ricevi questa email perché utilizzi il server Jellyfin di ${jellyfin_owner_name}. Se non desideri più ricevere queste email, puoi annullare l'iscrizione notificandolo a ${unsubscribe_email}.",
        "added_on": "Aggiunto il",
        "episodes": "Episodi",
        "episode": "Episodio",
        "new_episodes": "nuovi episodi",
        "footer_project_open_source": "è un progetto open source.",
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Sviluppato con ❤️ da <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> e <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">i collaboratori</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, con licenza AGPLv3.",
    },   
    "pt": {
        "discover_now": "Descubra agora",
        "new_film": "Novos filmes:",
        "new_tvs": "Novas séries:",
        "currently_available": "Atualmente disponível no Jellyfin:",
        "movies_label": "Filmes",
        "episodes_label": "Episódios",
        "footer_label": "Você está recebendo este e-mail porque está usando o servidor Jellyfin de ${jellyfin_owner_name}. Se quiser parar de receber esses e-mails, você pode cancelar a inscrição notificando ${unsubscribe_email}.",
        "added_on": "Adicionado em",
        "episodes": "Episódios",
        "episode": "Episódio",
        "new_episodes": "novos episódios",
        "footer_project_open_source": "é um projeto de código aberto.",
        # While the AGPLv3 license allows modification and redistribution, I kindly ask that the footer attribution remain intact to acknowledge the original project and its contributors. This helps support the open-source community and gives credit where it's due. Thanks !
        "footer_developed_by": 'Desenvolvido com ❤️ por <a href="https://github.com/SeaweedbrainCY/" class="footer-link">SeaweedbrainCY</a> e <a href="https://github.com/seaweedbraincy/jellyfin-newsletter/graphs/contributors" class="footer-link">colaboradores</a>.',
        "license_and_copyright": "Copyright © 2025 Nathan Stchepinsky, licenciado sob AGPLv3.",
    },
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
    with open(f"./themes/new_media/{configuration.conf.email_template.theme}/main.html", encoding='utf-8') as template_file:
        template = template_file.read()
        
        if configuration.conf.email_template.language in configuration.conf.email_template.available_lang:
            for key in TRANSLATIONS[configuration.conf.email_template.language]:
                template = re.sub(
                    r"\${" + key + "}", 
                    TRANSLATIONS[configuration.conf.email_template.language][key], 
                    template
                )
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