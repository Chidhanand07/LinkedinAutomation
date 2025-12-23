package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS processed_profiles (
			profile_url TEXT PRIMARY KEY,
			processed_at TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func IsProfileProcessed(db *sql.DB, profileURL string) bool {
	row := db.QueryRow(
		"SELECT COUNT(1) FROM processed_profiles WHERE profile_url = ?",
		profileURL,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

func MarkProfileProcessed(db *sql.DB, profileURL string) error {
	_, err := db.Exec(
		"INSERT OR IGNORE INTO processed_profiles(profile_url, processed_at) VALUES (?, ?)",
		profileURL,
		time.Now(),
	)
	return err
}
