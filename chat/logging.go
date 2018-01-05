package chat

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/treychua/beatricethetelegrambot/request"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Svc    Service
}

func (mw LoggingMiddleware) HandleRequest(r *request.Request) (s string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "HandleRequest",
			"input", fmt.Sprintf("r: %#v", r),
			"output", fmt.Sprintf("s: %v", s),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	// s, err = mw.Svc.HandleRequest(r)
	s, err = handleRequest(r, mw)
	return
}

func (mw LoggingMiddleware) getChat(r *request.Request) (c *chat, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "getChat",
			"input", fmt.Sprintf("r: %#v", r),
			"output", c,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	c, err = mw.Svc.getChat(r)
	return
}

func (mw LoggingMiddleware) handleAdd(c *chat, r *request.Request) (s string, err error) {
	defer func(begin time.Time) {
		mw.Logger.Log(
			"method", "handleAdd",
			"input", fmt.Sprintf("c: %#v, r: %#v", c, r),
			"output", s,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	s, err = mw.Svc.handleAdd(c, r)
	return
}
