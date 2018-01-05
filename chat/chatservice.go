package chat

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/treychua/beatricethetelegrambot/lunchvenue"

	"gopkg.in/mgo.v2/bson"

	"github.com/treychua/beatricethetelegrambot/request"
)

// Service handles any commands or inputs from the chat
type Service interface {
	HandleRequest(r *request.Request) (string, error)

	getChat(r *request.Request) (*chat, error)
	handleAdd(c *chat, r *request.Request) (string, error)
	// handleDelete(req *req.Request) (string, error)
	// handleList(req *req.Request) (string, error)
}

// ServiceImpl is an empty struct that'll have the methods for fulfilling Service's interface
type ServiceImpl struct{}

func (cs ServiceImpl) getChat(r *request.Request) (*chat, error) {
	return getChatFromDB(r)
}

func (cs ServiceImpl) handleAdd(c *chat, r *request.Request) (string, error) {
	if 2 > len(r.Message) {
		return "Sorry? I didn't catch what you said! Use '/add <your venue>'!", nil
	}
	location := strings.Join(r.Message[1:], " ")
	err := c.Venues.Add(location)
	if _, ok := err.(*lunchvenue.LocationAlreadyExistsError); ok {
		return `Location already exists~! Try somewhere else!  ¯\_(ツ)_/¯`, nil
	}
	err = updateChatTable(c, r)
	if nil != err {
		return "", err
	}

	return location + " added successfully~ Ehehe~", nil
}

// func (cs ServiceImpl) handleDelete(msg []string) (string, error) {
// }

// func (cs ServiceImpl) handleList(msg []string) (string, error) {
// }

// HandleRequest handles a request given by Telegram and returns a reply message.
func (cs ServiceImpl) HandleRequest(r *request.Request) (string, error) {
	return handleRequest(r, cs)
}

// =============================================================================
// helpers
// =============================================================================
func handleRequest(r *request.Request, s Service) (string, error) {
	if 0 == len(r.Message) {
		return "", errors.New("chatService::handleRequest() - message of length 0")
	}

	c, err := s.getChat(r)
	if nil != err {
		return "", errors.New("chatService::handleRequest() - chat cannot be retrieved")
	}

	var reply string
	switch r.Message[0] {

	case "/add":
		reply, err = s.handleAdd(c, r)
		if err != nil {
			return reply, err
		}

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
