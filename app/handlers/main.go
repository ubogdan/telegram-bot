package handlers

import (
	"log"
	"sync"

	telegram "gopkg.in/telegram-bot-api.v4"
)

var DefaultCommandHandler handler

type CommandHandlerFunc func(bot *telegram.BotAPI, msg *telegram.Message) error

type handler struct {
	sync.RWMutex
	commands map[string]CommandHandlerFunc
}

func (ch *handler) RegisterCommand(command string, handler CommandHandlerFunc) {
	ch.Lock()
	defer ch.Unlock()

	if ch.commands == nil {
		ch.commands = map[string]CommandHandlerFunc{
			"help": func(bot *telegram.BotAPI, msg *telegram.Message) error {
				_, err := bot.Send(telegram.NewMessage(msg.Chat.ID, "Type /settings."))
				return err
			},
		}
	}

	ch.commands[command] = handler
}

func (ch *handler) HandleMessage(bot *telegram.BotAPI, msg *telegram.Message) {
	ch.RLock()
	defer ch.RUnlock()

	if ch.commands == nil {
		log.Printf("no handler registered")
		return
	}

	commandName := msg.Command()

	command, found := ch.commands[commandName]
	if !found {
		_, err := bot.Send(telegram.NewMessage(msg.Chat.ID, "I don't know that command"))
		if err != nil {
			log.Println(err)
		}

		return
	}

	err := command(bot, msg)
	if err != nil {
		log.Println("error handling command", commandName, ":", err)
	}
}
