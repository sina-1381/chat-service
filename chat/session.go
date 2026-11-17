package chat

import (
	"chat-service/models"
	"chat-service/services"

	"github.com/gorilla/websocket"
)

// message structure
type Message struct {
	Msg    string      `json:"msg"`
	To     string      `json:"to"`
	From   interface{} `json:"from"`
	Type   string      `json:"type"`
	Status string
}

// websocket session structure
type Session struct {

	// list of users clients connected to the session
	Clients map[interface{}]*Client

	// broadcast channel for routing messages
	Broadcast chan *Message

	// client connecting to session
	Register chan *Client

	// client disconnecting from session
	UnRegister chan *Client
}

// websocket client structure
type Client struct {

	// session info in client side
	Hub *Session

	// connection to the session websocket
	Conn *websocket.Conn

	// message sending queue
	Send chan *Message

	// clients username data
	From   interface{}
	Status bool
}

func NewSession() *Session {
	return &Session{
		Clients:    make(map[interface{}]*Client),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

// the entire session side logic of websocket that is ran on a concurrent
// goroutine and its always running in background and only stops if
// project is stopped . this method is called in main file .
func (s *Session) Run() {
	for {
		select {

		// if client is registered to session this adds
		// the user client to clients list in session side
		case client := <-s.Register:
			s.Clients[client.From] = client
			client.Status = true

			// if user client unregisters remove from session
			// clients list and close the clients send channel
		case client := <-s.UnRegister:
			if _, ok := s.Clients[client.From]; ok {

				// delete from client list
				delete(s.Clients, client.From)

				// close send channel
				close(client.Send)
				client.Status = false
			}

			// broadcast channel in session open for messages
		case message := <-s.Broadcast:

			// if message type is privet
			if message.Type == "privet" {
				client := s.Clients[message.To]

				// check if reciver is connected to session if not drop message
				if client == nil {
					continue
				}
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.Clients, client.From)
				}

				// if message type is group:
			} else if message.Type == "group" {

				// get all users username in this groupchat from db and
				// send the message to their clients in session clients list
				var username []string
				services.CacheRun("users", username)
				if username == nil {
					models.DB.Table("group_chats").Select("users.user_name").
						Where("group_chats.title = ?", message.To).
						Joins("left join user_groups on group_chats.id = user_groups.group_chat_id").
						Joins("left join users on user_groups.user_id = users.id").Find(&username)
					services.CacheRun("users", username)
				}

				// range users in clients list. if user not connected to session drop message
				for _, e := range username {
					client := s.Clients[e]
					if client == nil {
						continue
					}
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(s.Clients, client.From)
					}
				}
			}
		}
	}
}
