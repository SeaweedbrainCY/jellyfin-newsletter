from source import configuration
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
        "added_on": "Added on"
    },
    "fr":{
        "discover_now": "Découvrir maintenant",
        "new_film": "Nouveaux films :",
        "new_tvs": "Nouvelles séries :",
        "currently_available": "Actuellement disponible sur Jellyfin :",
        "movies_label": "Films",
        "episodes_label": "Épisodes",
        "footer_label":"Vous recevez cet email car vous utilisez le serveur Jellyfin de ${jellyfin_owner_name}. Si vous ne souhaitez plus recevoir ces emails, vous pouvez vous désinscrire en notifiant ${unsubscribe_email}.",
        "added_on": "Ajouté le"
    }
}


def populate_email_template(movies, series, total_tv, total_movie) -> str:
    with open("./template/new_media_notification.html") as template:
        template=''.join(line.rstrip() for line in template)
        if configuration.conf.email_template.language in ["fr", "en"]:
            for key in translation[configuration.conf.email_template.language]:
                template = re.sub(r"\${"+key+"}", translation[configuration.conf.email_template.language][key], template)
        else:
            raise Exception(f"[FATAL] Language {configuration.conf.email_template.language} not supported. Supported languages are fr and en")

        custom_keys = [
            {"key":"title", "value":configuration.conf.email_template.title}, 
            {"key":"subtitle", "value":configuration.conf.email_template.subtitle},
            {"key":"jellyfin_url", "value":configuration.conf.email_template.jellyfin_url},
            {"key":"jellyfin_owner_name", "value":configuration.conf.email_template.jellyfin_owner_name},
            {"key":"unsubscribe_email", "value":configuration.conf.email_template.unsubscribe_email}
        ]
        for key in custom_keys:
            template = re.sub(r"\${"+key["key"]+"}", key["value"], template)

        if len(movies.keys())>=1:
            template = re.sub(r"\${display_movies}", "", template) 
            movies_html=""
            for movie_title in movies.keys():
                movies_html += f"""
                <div class="movie_container">
            <div class="movie_bg" style="background: url({movies[movie_title]["poster"]}) no-repeat; background-size: cover; background-position-y: center; ">
                <table class="movie" >
                <tr>
                <td class="img">
                    <img src="{movies[movie_title]["poster"]}" alt="Image du film" >
                </td>
                <td class="info">
                    <h3>{movie_title}</h3>
                    <span class="added">{translation[configuration.conf.email_template.language]["added_on"]} {movies[movie_title]["created_on"].split("T")[0]} </span>
                    <p class="description">{movies[movie_title]["description"]}</p>
                        
                    </td>
                    <tr>
                </table>
            </div>
        </div>
        """
            template = re.sub(r"\${films}", movies_html, template)

        else:
            template = re.sub(r"\${display_movies}", "display:none", template) 

        if len(series.keys())>=1:
            template = re.sub(r"\${display_tv}", "", template) 
            series_html=""
            for serie_title in series.keys():
                series_html += f"""
                <div class="movie_container">
            <div class="movie_bg" style="background: url({series[serie_title]["poster"]}) no-repeat; background-size: cover; background-position-y: center; ">
                <table class="movie" >
                <tr>
                <td class="img">
                    <img src="{series[serie_title]["poster"]}" alt="Image du film" >
                </td>
                <td class="info">
                    <h3>{serie_title} {", ".join(series[serie_title]["seasons"])}</h3>
                    <span class="added">{translation[configuration.conf.email_template.language]["added_on"]} {series[serie_title]["created_on"].split("T")[0]} </span>
                    <p class="description">{series[serie_title]["description"]}</p>
                        
                    </td>
                    </tr>
                </table>
            </div>
        </div>
        """
            template = re.sub(r"\${tvs}", series_html, template)
        else:
            template = re.sub(r"\${display_tv}", "display:none", template) 

        template = re.sub(r"\${series_count}", str(total_tv), template)
        template = re.sub(r"\${movies_count}", str(total_movie), template)
        return template




