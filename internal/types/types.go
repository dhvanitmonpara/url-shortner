package types

type URL struct {
	Id          int64  `json:"id"`
	OriginalURL string `json:"original_url" validate:"required"`
	ShortenURL  string `json:"shorten_url" validate:"required"`
	CreatedAt   int    `json:"created_at" validate:"required"`
}
