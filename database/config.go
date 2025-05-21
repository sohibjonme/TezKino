package database

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("Error for getting DATABASE_URL")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func BotToken() string {
	token := os.Getenv("BOT_TOKEN")
	return token
}
