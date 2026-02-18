package tmdb

type MediaType string

const (
	MediaTypeMovie  MediaType = "movie"
	MediaTypeSeries MediaType = "tv"
)

func (m MediaType) ToString() string {
	return string(m)
}
