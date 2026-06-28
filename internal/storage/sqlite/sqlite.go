package sqlite

import (
	"database/sql"
	"fmt"
	"url-shortner/internal/config"
	"url-shortner/internal/types"

	gonanoid "github.com/matoous/go-nanoid/v2"

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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (
  id TEXT PRIMARY KEY,
	redirect_to TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateURL(redirectTO string) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO urls (id, redirect_to) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 10)
	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(id, redirectTO)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetOriginalURLById(id string) (types.URL, error) {
	stmt, err := s.Db.Prepare("SELECT id, redirect_to, created_at FROM urls WHERE id = ? LIMIT 1")
	if err != nil {
		return types.URL{}, err
	}

	defer stmt.Close()

	var url types.URL

	fmt.Println("fetching url with id:", id)

	err = stmt.QueryRow(id).Scan(&url.Id, &url.RedirectTO, &url.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.URL{}, fmt.Errorf("no url found with id %s", fmt.Sprint(id))
		}
		return types.URL{}, fmt.Errorf("query error: %w", err)
	}

	return url, nil
}

func (s *Sqlite) GetURLs() ([]types.URL, error) {
	stmt, err := s.Db.Prepare("SELECT id, redirect_to, created_at FROM urls")
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

		err := rows.Scan(&url.Id, &url.RedirectTO, &url.CreatedAt)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}

func (s *Sqlite) UpdateUrl(id string, redirectTO string) (types.URL, error) {
	stmt, err := s.Db.Prepare("UPDATE urls SET redirect_to = ? WHERE id = ?")
	if err != nil {
		return types.URL{}, err
	}

	defer stmt.Close()

	var url types.URL

	err = stmt.QueryRow(redirectTO, id).Scan(&url.Id, &url.RedirectTO, &url.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.URL{}, fmt.Errorf("no url found with id %s", fmt.Sprint(id))
		}
		return types.URL{}, fmt.Errorf("query error: %w", err)
	}

	return url, nil
}
