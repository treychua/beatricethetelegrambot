package chat

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/treychua/beatricethetelegrambot/lunchvenue"
	"github.com/treychua/beatricethetelegrambot/request"
	"gopkg.in/mgo.v2/bson"
)

type chat struct {
	ChatID int64                  `bson:"chatid"`
	Venues lunchvenue.LunchVenues `bson:"venues"`
}

// GetChat retrieves an existing chat model from the database, and failing that, attempts to
// create a new one in the db
func GetChat(req *request.Request) *chat {
	s := req.Session.Copy()
	defer s.Close()

	collection := s.DB("beatricedb").C("chats")

	var ch chat
	err := collection.Find(bson.M{"chatid": req.ChatID}).One(&ch)
	if err != nil {
		errString := err.Error()
		if "not found" != errString {
			log.Println("err occurred", err)
			return nil
		}

		log.Printf("ChatID of %v not found, inserting new Chat", req.ChatID)
		ch = chat{req.ChatID, lunchvenue.LunchVenues{}}
		collection.Insert(&ch)
	}

	return &ch
}

func (c *chat) HandleRequest(r *request.Request) string {
	fmt.Println(r)

	if 0 == len(r.Message) {
		return ""
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
			log.Println(err)
			reply = err.Error()
		} else {
			updateChatTable(c, r)
			reply = location + " added successfully~"
		}

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

		reply = location + " removed successfully~\n"
		reply += "Remaining list of venues:\n"

		for i, v := range c.Venues {
			reply += strconv.Itoa(i+1) + ": " + v.Location + "\n"
		}

	case "get_random_lunch_venue":
		reply = "Feature not ready yet"
	}

	return reply
}

func updateChatTable(c *chat, r *request.Request) {
	s := r.Session.Copy()
	defer s.Close()

	collection := s.DB("beatricedb").C("chats")
	collection.Update(bson.M{"chatid": c.ChatID}, &c)
}
