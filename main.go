package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/treychua/beatricethetelegrambot/chat"
	"github.com/treychua/beatricethetelegrambot/request"
	mgo "gopkg.in/mgo.v2"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)

	var svc chat.ChatService
	svc = chat.ChatServiceImpl{}
	svc = chat.LoggingMiddleware{Logger: logger, Svc: svc}

	session, err := setupDB()
	logger.Log(
		"method", "setupDB",
		"output", fmt.Sprintf("%#v", session),
		"err", err,
	)
	defer session.Close()

	bot, updates, err := getTelegramUpdates()
	logger.Log(
		"method", "getTelegramUpdates",
		"output 1", fmt.Sprintf("%#v", bot),
		"output 2", fmt.Sprintf("%#v", updates),
		"err", err,
	)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messages := strings.Fields(update.Message.Text)
		newRequest := request.Request{Session: session, ChatID: chatID, Message: messages}

		reply, err := svc.HandleRequest(&newRequest)
		if nil == err {
			msg := tgbotapi.NewMessage(newRequest.ChatID, reply)
			bot.Send(msg)
		}

	}
}

// =============================================================================
// helper functions
// =============================================================================
func setupDB() (*mgo.Session, error) {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true) // not very sure what this does yet
	return session, nil
}

func getTelegramUpdates() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI("468474472:AAEoKhUM1ZTpNUSCOWExEsEXbwhkXLGapIg")
	if err != nil {
		return nil, nil, err
	}

	bot.Debug = false
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	return bot, updates, err
}
