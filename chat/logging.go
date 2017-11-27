package chat

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/treychua/beatricethetelegrambot/request"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Svc    ChatService
}

func (mw LoggingMiddleware) getChat(req *request.Request) (c *chat, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "getChat",
			"input", fmt.Sprintf("%#v", req),
			"output", c,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	c, err = mw.Svc.getChat(req)
	return
}

func (mw LoggingMiddleware) HandleRequest(r *request.Request) (s string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "HandleRequest",
			"input", fmt.Sprintf("%#v", r),
			"output", s,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	s, err = mw.Svc.HandleRequest(r)
	return
}
