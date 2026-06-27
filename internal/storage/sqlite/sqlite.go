package sqlite

import (
	"database/sql"
	"fmt"
	"time"
	"url-shortner/internal/config"
	"url-shortner/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateURL(originalUrl string, shortenUrl string) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO urls (original_url, shorten_url, created_at) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(originalUrl, shortenUrl, time.UTC)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetOriginalURLById(id int64) (types.URL, error) {
	stmt, err := s.Db.Prepare("SELECT id, original_url, shorten_url, created_at FROM urls WHERE id = ? LIMIT 1")
	if err != nil {
		return types.URL{}, err
	}

	defer stmt.Close()

	var url types.URL

	err = stmt.QueryRow(id).Scan(&url.Id, &url.OriginalURL, &url.ShortenURL, &url.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.URL{}, fmt.Errorf("no url found with id %s", fmt.Sprint(id))
		}
		return types.URL{}, fmt.Errorf("query error: %w", err)
	}

	return url, nil
}

func (s *Sqlite) GetOriginalURLByShortenURL(shortenUrl string) (types.URL, error) {
	stmt, err := s.Db.Prepare("SELECT id, original_url, shorten_url, created_at FROM urls WHERE shorten_url = ? LIMIT 1")
	if err != nil {
		return types.URL{}, err
	}

	defer stmt.Close()

	var url types.URL

	err = stmt.QueryRow(shortenUrl).Scan(&url.Id, &url.OriginalURL, &url.ShortenURL, &url.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.URL{}, fmt.Errorf("no url found with shorten url %s", shortenUrl)
		}
		return types.URL{}, fmt.Errorf("query error: %w", err)
	}

	return url, nil
}

func (s *Sqlite) GetURLs() ([]types.URL, error) {
	stmt, err := s.Db.Prepare("SELECT id, original_url, shorten_url, created_at FROM urls")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var urls []types.URL

	for rows.Next() {
		var url types.URL

		err := rows.Scan(&url.Id, &url.OriginalURL, &url.ShortenURL, &url.CreatedAt)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}
