package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/ubogdan/telegram-bot/app"
	telegram "gopkg.in/telegram-bot-api.v4"
)

func main() {
	botToken, found := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !found {
		log.Fatalf("BOT_TOKEN not found")
	}

	bot, err := telegram.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("newBot %s", err)
	}

	updates := bot.ListenForWebhook("/")

	go lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)

	for update := range updates {
		go app.HandleMessage(bot, update.Message)
	}

}
