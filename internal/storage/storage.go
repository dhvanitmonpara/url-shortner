package storage

import "url-shortner/internal/types"

type Storage interface {
	CreateURL(originalUrl string, shortenUrl string) (int64, error)
	GetOriginalURLById(id int64) (types.URL, error)
	GetOriginalURLByShortenURL(shortenUrl string) (types.URL, error)
	GetURLs() ([]types.URL, error)
}
