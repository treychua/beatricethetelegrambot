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
	gopath := os.Getenv("HOME")
	logpath := gopath + "/logs.txt"

	var f *os.File
	var err error

	if f, err = os.Open(logpath); err != nil {
		f, err = os.Create(logpath)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()

	logger := log.NewLogfmtLogger(f)

	// fieldKeys := []string{"method", "error"}
	// requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	// 	Namespace: "chat",
	// 	Subsystem: "chat_service",
	// 	Name:      "request_count",
	// 	Help:      "Number of requests received.",
	// }, fieldKeys)
	// requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	// 	Namespace: "chat",
	// 	Subsystem: "chat_service",
	// 	Name:      "request_latency_microseconds",
	// 	Help:      "Total duration of requests in microseconds.",
	// }, fieldKeys)
	// countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	// 	Namespace: "chat",
	// 	Subsystem: "chat_service",
	// 	Name:      "count_result",
	// 	Help:      "The result of each count method.",
	// }, []string{})

	var svc chat.Service
	svc = chat.ServiceImpl{}
	svc = chat.LoggingMiddleware{Logger: logger, Svc: svc}
	// svc = chat.InstrumentingMiddleware{
	// 	RequestCount:   requestCount,
	// 	RequestLatency: requestLatency,
	// 	CountResult:    countResult,
	// 	Svc:            svc,
	// }

	session, err := setupDB()
	logger.Log(
		"method", "setupDB",
		"output", fmt.Sprintf("%#v", session),
		"err", err,
	)
	defer session.Close()

	bot, updates, err := getTelegramUpdates()
	lastid := 0

	logger.Log(
		"method", "getTelegramUpdates",
		"output 1", fmt.Sprintf("%#v", bot),
		"output 2", fmt.Sprintf("%#v", updates),
		"err", err,
	)

	for update := range updates {

		if update.Message == nil && (update.CallbackQuery == nil || lastid == 0) {
			continue
		}

		var chatID int64
		var messages []string
		var reply string
		var err error
		msg := tgbotapi.NewMessage(0, "")

		if lastid != 0 && update.CallbackQuery != nil {
			//choose next random location
			if update.CallbackQuery.Data == "next" {

				chatID = update.CallbackQuery.Message.Chat.ID
				messages = strings.Fields("/random")
				newRequest := request.Request{Session: session, ChatID: chatID, Message: messages}
				reply, err = svc.HandleRequest(&newRequest)
				msg = tgbotapi.NewMessage(newRequest.ChatID, reply)
				butt := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Okay~", "ok"), tgbotapi.NewInlineKeyboardButtonData("Choose again!", "next"))
				keyb := tgbotapi.NewInlineKeyboardMarkup(butt)
				// keyb.OneTimeKeyboard = true
				msg.ReplyMarkup = &keyb

			} else {
				//save something to db?
				reply = "Yay! Enjoy your lunch~~~"
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
			}

		} else {
			//update.Message != nil
			chatID = update.Message.Chat.ID
			messages = strings.Fields(update.Message.Text)
			newRequest := request.Request{Session: session, ChatID: chatID, Message: messages}
			reply, err = svc.HandleRequest(&newRequest)
			msg = tgbotapi.NewMessage(newRequest.ChatID, reply)

			if len(messages) > 0 && strings.Contains(strings.ToLower(messages[0]), "/random") {
				butt := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Okay", "ok"), tgbotapi.NewInlineKeyboardButtonData("Choose again!", "next"))
				keyb := tgbotapi.NewInlineKeyboardMarkup(butt)
				msg.ReplyMarkup = &keyb
			}
		}

		if nil == err {
			sm, _ := bot.Send(msg)
			lastid = sm.MessageID
		}

		// logger.Log(
		// 	"request_count", fmt.Sprintf("%#v", requestCount),
		// 	"request_latency", fmt.Sprintf("%#v", requestLatency),
		// 	"count_result", fmt.Sprintf("%#v", countResult),
		// )

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
