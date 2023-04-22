//go:build lambda
// +build lambda

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ubogdan/telegram-bot/app/handlers"
	"github.com/ubogdan/telegram-bot/app/httpadapter"
	telegram "gopkg.in/telegram-bot-api.v4"
)

func main() {
	botToken, found := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !found {
		log.Fatalf("BOT_TOKEN not found")
	}

	botHookURL, found := os.LookupEnv("")
	if !found {
		log.Fatalf("TELEGRAM_BOT_WEBHOOK not found")
	}

	bot, err := telegram.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("newBot %s", err)
	}

	// Register http handler to http.DefaultServeMux
	updates := bot.ListenForWebhook("/")

	// Register hook with telegram
	_, err = bot.SetWebhook(telegram.NewWebhook(botHookURL))
	if err != nil {
		panic(err)
	}

	// Start listening for updates
	go func() {
		for update := range updates {
			go handlers.DefaultCommandHandler.HandleMessage(bot, update.Message)
		}
	}()

	lambda.Start(httpadapter.New(http.DefaultServeMux))
}
