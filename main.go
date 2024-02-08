package main

import (
	"discord-bot/bot"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file.", err)
		return
	}

	bot.Token = os.Getenv("TOKEN")
	bot.ServerID = os.Getenv("SERVER_ID")
	bot.TargetUsers = []string{os.Getenv("TARGET_USER1"), os.Getenv("TARGET_USER2")}

	bot.Run() // call the run function of bot/bot.go
}
