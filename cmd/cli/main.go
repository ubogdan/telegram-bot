package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ubogdan/telegram-bot/app/handlers"
	"golang.org/x/crypto/acme/autocert"
	telegram "gopkg.in/telegram-bot-api.v4"
)

func main() {
	botToken, found := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !found {
		log.Fatalf("TELEGRAM_BOT_TOKEN not found")
	}

	webHookDomain, found := os.LookupEnv("TELEGRAM_BOT_WEBHOOK_DOMAIN")
	if !found {
		log.Fatalf("TELEGRAM_BOT_WEBHOOK_DOMAIN not found")
	}

	bot, err := telegram.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("newBot %s", err)
	}

	log.Println("Authorized on account", bot.Self.UserName)

	_, err = bot.SetWebhook(telegram.NewWebhook("https://" + webHookDomain + "/"))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	dataDir := "."
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		//Email:  "ubogdan@gmail.com",
		HostPolicy: func(ctx context.Context, host string) error {
			// Note: change to your real host
			allowedHost := webHookDomain
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		},
		Cache: autocert.DirCache(dataDir),
	}

	httpsSrv := http.Server{
		Addr:         ":443",
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate, ServerName: webHookDomain, NextProtos: []string{"h2", "http/1.1"}},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      http.DefaultServeMux,
	}
	go func() {
		fmt.Printf("Starting HTTPs server on %s\n", httpsSrv.Addr)
		err = httpsSrv.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("failed to start https server: %s", err)
		}
	}()

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)

	httpSrv := http.Server{
		Addr:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      m.HTTPHandler(mux),
	}

	go func() {
		fmt.Printf("Starting HTTP server on %s\n", httpSrv.Addr)
		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	}()

	for update := range updates {
		go handlers.DefaultCommandHandler.HandleMessage(bot, update.Message)
	}

}
