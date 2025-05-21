package main

import (
	"log"
	"tezkinobot/database"
	"tezkinobot/handler"
	"tezkinobot/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env fileni o'qishdagi xatolik")
	}
	db, err := database.InitDB()
	defer db.Close()
	if err != nil {
		log.Println("NImadir")
		log.Fatal(err)
	}
	log.Println(db)
	server.CreateDatabase(db)
	handler.LogicBot(db)

}
