package chat

import (
	"errors"
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

	case "/add_lunch_venue":
		if 2 > len(r.Message) {
			break
		}
		location := strings.Join(r.Message[1:], " ")
		err := c.Venues.InsertLunchVenue(location)
		if nil != err {
			return reply, err
		}
		err = updateChatTable(c, r)
		if nil != err {
			return reply, err
		}

		reply = location + " added successfully~"

	case "/list_lunch_venues":
		reply = "List of venues:\n"

		for i, v := range c.Venues {
			reply += strconv.Itoa(i+1) + ": " + v.Location + "\n"
		}

	case "/delete_lunch_venue":
		if 2 > len(r.Message) {
			break
		}

		location := strings.Join(r.Message[1:], " ")
		c.Venues.DeleteLunchVenue(location)

		updateChatTable(c, r)
		if nil != err {
			return reply, err
		}

		reply = location + " removed successfully~\n"
		reply += "Remaining list of venues:\n"

		for i, v := range c.Venues {
			reply += strconv.Itoa(i+1) + ": " + v.Location + "\n"
		}

	case "get_random_lunch_venue":
		reply = "Feature not ready yet"
	}

	return reply, nil
}

func updateChatTable(c *chat, r *request.Request) error {
	s := r.Session.Copy()
	defer s.Close()

	collection := s.DB("beatricedb").C("chats")
	return collection.Update(bson.M{"chatid": c.ChatID}, &c)
}
