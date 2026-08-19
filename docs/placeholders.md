# Placeholders

You may want to use dynamic values in your newsletter — that's what placeholders are for. For every custom string in `config.yml`, you can insert a placeholder that gets replaced with a real value in the final email.

## Available placeholders

| Placeholder | Old placeholder (`< v1.0.0`) | Example | Description | Version |
|---|---|---|---|---|
| `{{.Date}}` | `{date}` | `2025-06-19` | Date of the day. Format `Y-m-d`. | `>= v1.0.0` |
| `{{.DayName}}` | `{day_name}` | `Monday` | Today's day name, localized. | `>= v1.0.0` |
| `{{.DayNumber}}` | `{day_number}` | `19` | Today's day number of the month. | `>= v1.0.0` |
| `{{.MonthName}}` | `{month_name}` | `June` | Current month name, localized. | `>= v1.0.0` |
| `{{.MonthNumber}}` | `{month_number}` | `06` | Current month number. | `>= v1.0.0` |
| `{{.Year}}` | `{year}` | `2025` | Current year. | `>= v1.0.0` |
| `{{.StartDate}}` | `{start_date}` | `2025-05-18` | First date new additions are taken into account. Depends on the `observed_period_days` config parameter. Format `Y-m-d`. | `>= v1.0.0` |
| `{{.StartDayName}}` | `{start_day_name}` | `Sunday` | Day name of the first observed date. | `>= v1.0.0` |
| `{{.StartDayNumber}}` | `{start_day_number}` | `19` | Day number of the first observed date. | `>= v1.0.0` |
| `{{.StartMonthName}}` | `{start_month_name}` | `May` | Month name of the first observed date. | `>= v1.0.0` |
| `{{.StartMonthNumber}}` | `{start_month_number}` | `05` | Month number of the first observed date. | `>= v1.0.0` |
| `{{.StartYear}}` | `{start_year}` | `2025` | Year of the first observed date. | `>= v1.0.0` |

?> Don't see what you're looking for? [Open an issue](https://github.com/SeaweedbrainCY/jellyfin-newsletter/issues) to request a new placeholder.

!> If you use a placeholder that doesn't exist, its key (e.g. `{{.ExamplePlaceholder}}`) will **not** be replaced and will show up literally in the final email.

## Example

Let's customize the email subject. In `config.yml`:

```yaml
[...]
email_template:
  language: "en"
  subject: "New addition of {{.MonthName}}"
[...]
```

Result:

<p align="center">
  <img width="261" alt="Rendered subject line showing the month name placeholder replaced" src="https://github.com/user-attachments/assets/33a8909d-fd00-4a6f-8292-89f3a7618da5">
</p>

?> Change the `language` field to one of the [supported languages](https://github.com/SeaweedbrainCY/jellyfin-newsletter#supported-languages) to get localized day and month names.

## See also

- [Configuration parameters](configuration.md#email_template) — where these placeholders can be used
- [Local theme files guide](local-themes.md) — placeholders are also available inside custom theme HTML
