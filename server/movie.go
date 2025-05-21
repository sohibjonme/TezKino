package server

import (
	"database/sql"
	"log"
	"strings"
)

func CreateDatabase(db *sql.DB) {
	query :=
		`CREATE TABLE IF NOT EXISTS movies(
					id SERIAL PRIMARY KEY,
					title TEXT NOT NULL,
					file_id TEXT NOT NULL
				)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertMovie(db *sql.DB, caption, fileID string) {
	query := `INSERT INTO movies (title, file_id) VALUES ($1, $2)`
	_, err := db.Exec(query, strings.ToLower(caption), fileID)
	if err != nil {
		log.Fatal("Error inserting informations")
	}
}

func SearchMovie(db *sql.DB, title string) (string, string) {
	pattern := "%" + title + "%"
	query := `
				SELECT file_id, title FROM movies WHERE LOWER(title) LIKE LOWER($1) LIMIT 1
			`
	row := db.QueryRow(query, pattern)
	var fileID string
	var caption string
	err := row.Scan(&fileID, &caption)
	if err != nil {
		return "", ""
	}
	return fileID, caption
}
