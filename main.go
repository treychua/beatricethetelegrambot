package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/treychua/beatricethetelegrambot/chat"
	"github.com/treychua/beatricethetelegrambot/request"
	mgo "gopkg.in/mgo.v2"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	var svc chat.ChatService
	svc = chat.ChatServiceImpl{}
	svc = chat.LoggingMiddleware{Logger: logger, Svc: svc}
	svc = chat.InstrumentingMiddleware{
		RequestCount:   requestCount,
		RequestLatency: requestLatency,
		CountResult:    countResult,
		Svc:            svc,
	}

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
