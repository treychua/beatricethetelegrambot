package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/treychua/beatricethetelegrambot/chat"
	"github.com/treychua/beatricethetelegrambot/request"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true) // not very sure what this does yet

	bot, err := tgbotapi.NewBotAPI("468474472:AAEoKhUM1ZTpNUSCOWExEsEXbwhkXLGapIg")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messages := strings.Fields(update.Message.Text)
		newRequest := request.Request{Session: session, ChatID: chatID, Message: messages}
		c := chat.GetChat(&newRequest)

		if nil != err {
			continue
		}

		reply := c.HandleRequest(&newRequest)
		if 0 != len(reply) {
			msg := tgbotapi.NewMessage(newRequest.ChatID, reply)
			bot.Send(msg)
		}

	}
}
