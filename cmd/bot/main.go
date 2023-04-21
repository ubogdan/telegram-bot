package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
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

func handleMessage(bot *telegram.BotAPI, msg *telegram.Message) {
	log.Println("update", msg.From, msg.From.ID, msg.From.LastName)
	if msg == nil { // ignore any non-Message Updates
		return
	}

	if !msg.IsCommand() { // ignore any non-command Messages
		return
	}

	_, err := bot.Send(telegram.NewMessage(msg.Chat.ID, handleCommand(bot, msg)))
	if err != nil {
		panic(err)
	}
}

func main() {
	botToken, found := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !found {
		log.Fatalf("BOT_TOKEN not found")
	}

	bot, err := telegram.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	go lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)

	for update := range bot.ListenForWebhook("/") {
		go handleMessage(bot, update.Message)
	}

}
