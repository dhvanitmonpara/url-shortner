package types

type URL struct {
	Id         string `json:"id,omitempty"`
	RedirectTO string `json:"redirect_to" validate:"required"`
	CreatedAt  string `json:"created_at,omitempty"`
}
