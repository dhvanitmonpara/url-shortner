package storage

import "url-shortner/internal/types"

type Storage interface {
	CreateURL(redirectTO string) (int64, error)
	GetOriginalURLById(id string) (types.URL, error)
	GetURLs() ([]types.URL, error)
	UpdateUrl(id string, redirectTo string) (types.URL, error)
}
