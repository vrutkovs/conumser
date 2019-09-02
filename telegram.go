package main

import (
	"log"
	"os"
	"strconv"

	tg "gopkg.in/telegram-bot-api.v4"
)

type tgBot struct {
	bot *tg.BotAPI
}

func createBot() (*tgBot, error) {
	botToken := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &tgBot{bot: bot}, nil
}

func (t *tgBot) sendMessage(room, message string) error {
	roomID, err := strconv.Atoi(room)
	if err != nil {
		return err
	}
	tgMessage := tg.NewMessage(int64(roomID), message)
	tgMessage.ParseMode = "Markdown"
	_, err = t.bot.Send(tgMessage)
	return err
}
