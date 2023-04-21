package app

import (
	"log"

	telegram "gopkg.in/telegram-bot-api.v4"
)

func handleCommand(bot *telegram.BotAPI, msg *telegram.Message) string {
	log.Println("handleCommand", msg.Command(), msg.CommandArguments())
	switch msg.Command() {
	case "help":
		return "type /settings."
	case "settings":
		return "I know nothing about settings"
	default:
		return "I don't know that command"
	}
}

func HandleMessage(bot *telegram.BotAPI, msg *telegram.Message) {
	if msg == nil { // ignore any non-Message Updates
		return
	}

	log.Println("update", msg.From, msg.From.ID, msg.From.LastName)

	if !msg.IsCommand() { // ignore any non-command Messages
		return
	}

	_, err := bot.Send(telegram.NewMessage(msg.Chat.ID, handleCommand(bot, msg)))
	if err != nil {
		panic(err)
	}
}
