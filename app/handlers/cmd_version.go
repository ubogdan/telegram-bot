package handlers

import (
	"log"

	telegram "gopkg.in/telegram-bot-api.v4"
)

func init() {
	DefaultCommandHandler.RegisterCommand("version", versionCMD)
}

func versionCMD(bot *telegram.BotAPI, msg *telegram.Message) error {
	_, err := bot.Send(telegram.NewMessage(msg.Chat.ID, "0.0.1"))
	if err != nil {
		log.Printf("bot.Send error %s", err)
	}

	return nil
}
