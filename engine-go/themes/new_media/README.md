# New media themes

This folder contains the email templates for the new media notification feature.

## Email template themes 

## Contribution

The template theme and organization are free to use and modify. Feel free to contribute by creating a pull request to create new themes or modify existing ones.

HTML templates can omit information, but cannot add new information. For example, a template can choose to not display the media description, but cannot add a new section for media ratings.

### Create a new theme

> [!IMPORTANT]
> Template theme and format are truly free. **But please, ensure to respect the global structure and the global features. The goal is to provide a variety of templates while keeping the same features and information.** If you want to add or configure new features, please open an issue or a discussion to talk about it.

> [!WARNING]
> Designing HTML email is tricky and not as straightforward as designing a web page. Email clients have many limitations and quirks. Please test your template in multiple email clients to ensure compatibility. Kindly review email clients documentation for limitations.

Templates should be designed using Go text/template syntax. This allows the script to replace the placeholders with the actual values when generating the email. The available placeholders are listed below.

Templates must have the following structure:
- A folder named after the template (e.g. `modern`, `classic`, etc.)
- Inside the folder, an `html`file named after the theme name (e.g. `modern.html`, `classic.html`, etc.) containing the HTML code and implementing (or not) the following placeholders: 

    - **Global Configuration**
        - `{{.HTMLLang}}` - HTML language attribute (e.g., "en", "fr")
        - `{{.HTMLDir}}` - HTML text direction ("ltr" or "rtl")

    - **Header Section**
        - `{{.Title}}` - Main email title/heading
        - `{{.Subtitle}}` - Subtitle text below the main title
        - `{{.JellyfinURL}}` - URL to the Jellyfin instance
        - `{{.DiscoverNowLabel}}` - Text label for the CTA button

    - **Movies Section**
        - `{{.DisplayNewMovies}}` - Boolean to show/hide movies section
        - `{{.NewFilmLabel}}` - Section heading for new movies
        - `{{.NewMovies}}` - Array of movie objects with:
            - `{{.Name}}` - Movie title
            - `{{.PosterURL}}` - Movie poster image URL
            - `{{.AddedOnLabel}}` - "Added on" text label
            - `{{.AdditionDate}}` - Date the movie was added
            - `{{.Overview}}` - Movie synopsis/description
            - `{{.IncludeItemOverviews}}` - Boolean to show/hide overview text

    - **TV Series Section**
        - `{{.DisplayNewSeries}}` - Boolean to show/hide TV series section
        - `{{.NewSeriesLabel}}` - Section heading for new series
        - `{{.NewSeries}}` - Array of series objects with:
            - `{{.SeriesName}}` - Series title (for alt text)
            - `{{.NewSeriesTitle}}` - Series display title
            - `{{.PosterURL}}` - Series poster image URL
            - `{{.AddedOnLabel}}` - "Added on" text label
            - `{{.AdditionDate}}` - Date the series was added
            - `{{.Overview}}` - Series synopsis/description
            - `{{.IncludeItemOverviews}}` - Boolean to show/hide overview text

    - **Statistics Section**
        - `{{.CurrentlyAvailableLabel}}` - Title for stats section
        - `{{.MoviesCount}}` - Total number of movies available
        - `{{.MoviesLabel}}` - Label for movies count
        - `{{.SeriesCount}}` - Total number of series available
        - `{{.SeriesLabel}}` - Label for series count

    - **Footer Section**
        - `{{.FooterLabel}}` - Main footer text
        - `{{.FooterProjectLinkLabel}}` - Link text for project repository
        - `{{.FooterOpenSourceProjectLabel}}` - Open source project attribution text
        - `{{.FooterDevelopedByLabel}}` - Developer attribution text
        - `{{.FooterLicenceAndCopyright}}` - License and copyright information


> [!IMPORTANT]
> It would be appreciated to include in your template footer the name and/or a link towards this repository. Open source projects thrive on visibility and contributions. Thank you!