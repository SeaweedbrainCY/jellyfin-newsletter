# New media themes

This folder contains the email templates for the new media notification feature.

## Email template themes 

## Contribution

The template theme and organization are free to use and modify. Feel free to contribute by creating a pull request to create new themes or modify existing ones.

HTML templates can omit information, but cannot add new information. For example, a template can choose to not display the media description, but cannot add a new section for media ratings.

### Create a new theme

> [!IMPORTANT]
> Template theme and format are truly free. **But please, ensure to respect the global structure and the gloal features. The goal is to provide a variety of templates while keeping the same features and information.** If you want to add or configure new features, please open an issue or a discussion to talk about it.

> [!WARNING]
> Designing HTML email is tricky and not as straightforward as designing a web page. Email clients have many limitations and quirks. Please test your template in multiple email clients to ensure compatibility. Kindly review email clients documentation for limitations.

Templates must have the following structure:
- A folder named after the template (e.g. `modern`, `classic`, etc.)
- Inside the folder, create a `main.html` file for the HTML version of the email. This file will be the structure for the email file. It must contain the CSS styles and the HTML structure. The available placeholders are : 
    - `${title}`: The title of the newsletter. Defined in the config file. 
    - `${subtitle}`: The subtitle of the newsletter. Defined in the config file. 
    - `${jellyfin_url}`: The URL of the Jellyfin server. Defined in the config file.
    - `${display_movies}`: Must be used in the `style` attribute of the movie section container. It will be replaced with `display:block;` or `display:none;` depending on whether there are new movies or not.
    - `${new_film}`: The translated title for the new movies section. Automatically translated by the script.
    - `${films}`: **Warning. This is a special placeholder**. This placeholder will be replaced by the HTML code for ALL new movies. The HTML structure for ONE movie is defined in the `movie.html` file. 
    - `${display_tv}`: Must be used in the `style` attribute of the series section container. It will be replaced with `display:block` or `display:none` depending on whether there are new series or not.
    - `${new_tvs}`: The translated title for the new series section. Automatically translated by the script.
    - `${tvs}`: **Warning. This is a special placeholder**. This placeholder will be replaced by the HTML code for ALL new series. The HTML structure for ONE series is defined in the `tv.html` file.
    - `${movies_count}`: The number of movies available in Jellyfin. Fetched from Jellyfin.
    - `${movies_label}`: Translated label for "movies". Automatically translated by the script.
    - `${series_count}`: The number of series available in Jellyfin. Fetched from Jellyfin.
    - `${episodes_label}`: Translated label for "episodes". Automatically translated by the script.
    - `${footer_label}`: Legal mention text. **Mandatory placeholder**. Defined in the config file.
- Create a `movie.html` file for the HTML structure of a single movie. This file will be used to generate the HTML code for each new movie. The final computed HTML code, concatenation of every `movie.html` filled structure, will be inserted in place of the `${films}` placeholder in the `main.html` file. It can contain the following placeholders:
    - `${movie_name}`: The name of the movie.
    - `${movie_overview}`: The description of the movie.
    - `${movie_overview_style}`: Must be used in the `style` attribute of the description container. It will be replaced with `display:block;` or `display:none;` depending on whether the description should be included or not (this is decided by the user in the config file).
    - `${movie_added_on}`: The release date of the movie.
    - `${movie_added_on_label}`: The translated label for "Release date". Automatically translated by the script.
    - `${movie_poster}`: The TMDB URL of the movie poster.

- Create a `tv.html` file for the HTML structure of a single series. This file will be used to generate the HTML code for each new series. The final computed HTML code, concatenation of every `tv.html` filled structure, will be inserted in place of the `${tvs}` placeholder in the `main.html` file. It can contain the following placeholders:
    - `${tv_title}`: The name of item, with the series name and the season/episode number (e.g. "Breaking Bad - Season 1, 2 & 3").
    - `${tv_overview}`: The description of the series.
    - `${tv_overview_style}`: Must be used in the `style` attribute of the description container. It will be replaced with `display:block;` or `display:none;` depending on whether the description should be included or not (this is decided by the user in the config file).
    - `${tv_added_on}`: The release date of the series.
    - `${tv_added_on_label}`: The translated label for "Release date". Automatically translated by the script.
    - `${tv_poster}`: The TMDB URL of the series poster.




> [!IMPORTANT]
> It would be appreciated to include in your template footer the name and/or a link towards this repository. Open source projects thrive on visibility and contributions. Thank you!