package chat

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/treychua/beatricethetelegrambot/request"
	"gopkg.in/mgo.v2/bson"
)

type ChatService interface {
	getChat(req *request.Request) (*chat, error)
	HandleRequest(r *request.Request) (string, error)
}

type ChatServiceImpl struct{}

func (cs ChatServiceImpl) getChat(req *request.Request) (*chat, error) {
	return getChatFromDB(req)
}

func (cs ChatServiceImpl) HandleRequest(r *request.Request) (string, error) {
	if 0 == len(r.Message) {
		return "", errors.New("chatService::HandleRequest() - message of length 0")
	}

	c, err := cs.getChat(r)
	if nil != err {
		return "", errors.New("chatService::HandleRequest() - chat cannot be retrieved")
	}

	var reply string
	switch r.Message[0] {

	case "/add":
		if 2 > len(r.Message) {
			break
		}
		location := strings.Join(r.Message[1:], " ")
		err := c.Venues.Add(location)
		if nil != err {
			return reply, err
		}
		err = updateChatTable(c, r)
		if nil != err {
			return reply, err
		}

		reply = location + " added successfully~"

	case "/list":
		reply = "List of venues:\n"

		for i, v := range c.Venues {
			reply += strconv.Itoa(i+1) + ": " + v.Location + "\n"
		}

	case "/remove":
		fallthrough

	case "/delete":
		if 2 > len(r.Message) {
			break
		}

		location := strings.Join(r.Message[1:], " ")
		c.Venues.Delete(location)

		err = updateChatTable(c, r)
		if nil != err {
			return reply, err
		}

		reply = location + " removed successfully~\n"
		reply += "Remaining list of venues:\n"

		for i, v := range c.Venues {
			reply += strconv.Itoa(i+1) + ": " + v.Location + "\n"
		}

	case "/random":

		i := rand.Intn(len(c.Venues))
		reply = c.Venues[i].Location
	}

	return reply, nil
}

func updateChatTable(c *chat, r *request.Request) error {
	s := r.Session.Copy()
	defer s.Close()

	collection := s.DB("beatricedb").C("chats")
	return collection.Update(bson.M{"chatid": c.ChatID}, &c)
}
