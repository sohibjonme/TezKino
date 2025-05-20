package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env fileni o'qishdagi xatolik")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set to .env file")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	initDB(db)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("Tokenni olishdagi xatolik")
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.ChannelPost != nil && update.ChannelPost.Video != nil {
			log.Println("dmskasssssssssssssssss")
			caption := update.ChannelPost.Caption
			fileID := update.ChannelPost.Video.FileID
			if caption != "" && fileID != "" {
				insert(db, caption, fileID)
			}
			continue
		}
		if update.Message == nil {
			continue
		}
		text := strings.ToLower(update.Message.Text)
		if text == "/start" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Menga kino nomini jo'nat. Men esa senga kinoni jo'nataman"))
			continue
		}
		fileID := searchMovie(db, text)
		if fileID == "" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Bu kino topilmadi"))
		} else {
			msg := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileID(fileID))
			bot.Send(msg)
		}

	}
}

func searchMovie(db *sql.DB, title string) string {
	pattern := "%" + title + "%"
	query := `
				SELECT file_id FROM movies WHERE LOWER(title) LIKE LOWER($1) LIMIT 1
			`
	row := db.QueryRow(query, pattern)
	var fileID string
	err := row.Scan(&fileID)
	if err != nil {
		return ""
	}
	return fileID
}

func insert(db *sql.DB, caption, fileID string) {
	query := `INSERT INTO movies (title, file_id) VALUES ($1, $2)`
	_, err := db.Exec(query, strings.ToLower(caption), fileID)
	if err != nil {
		log.Fatal("Error inserting informations")
	}
}

func initDB(db *sql.DB) {
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
