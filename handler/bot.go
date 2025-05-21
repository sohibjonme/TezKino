package handler

import (
	"database/sql"
	"log"
	"strings"
	"tezkinobot/database"
	"tezkinobot/server"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LogicBot(db *sql.DB) {
	token := database.BotToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Tokenni olishdagi xatolik")
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		desc := tgbotapi.SetChatDescriptionConfig{
			ChatID:      update.Message.Chat.ID,
			Description: "nimadir"}
		bot.Request(desc)
		if update.ChannelPost != nil && update.ChannelPost.Video != nil {

			caption := update.ChannelPost.Caption
			fileID := update.ChannelPost.Video.FileID
			if caption != "" && fileID != "" {
				server.InsertMovie(db, caption, fileID)
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
		fileID, caption := server.SearchMovie(db, text)
		if fileID == "" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Bu kino topilmadi"))
		} else {
			msg := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileID(fileID))
			msg.Caption = caption

			sentMsg, err := bot.Send(msg)
			alter := tgbotapi.NewMessage(update.Message.Chat.ID, "Bu kino 1 soatda o'chiriladi. Tezda yuklab oling!!!")
			bot.Send(alter)
			if err != nil {
				log.Fatal("Error sending video")
			}
			go func(chatID int64, messageID int) {
				time.Sleep(1 * time.Hour)
				delConf := tgbotapi.DeleteMessageConfig{
					ChatID:    chatID,
					MessageID: messageID,
				}
				if _, err := bot.Request(delConf); err != nil {
					log.Println("failed to delete message", err)
				} else {
					log.Println("Deleted Video")
				}

			}(sentMsg.Chat.ID, sentMsg.MessageID)
		}

	}
}
