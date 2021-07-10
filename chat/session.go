package chat

import (
	"ginGorm/models"
	"ginGorm/services"
	"github.com/gorilla/websocket"
)

type Message struct {
	Msg string `json:"msg"`
	To string `json:"to"`
	From interface{} `json:"from"`
	Type string `json:"type"`
	Status string
}

type Session struct {
	Clients map[interface{}]*Client
	Broadcast chan *Message
	Register chan *Client
	UnRegister chan *Client
}

type Client struct {
	Hub *Session
	Conn *websocket.Conn
	Send chan *Message
	From interface{}
	Status bool
}

func NewSession () *Session{
	return &Session{
		Clients:    make(map[interface{}]*Client),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

func (s *Session)Run() {
	for {
		select {
		case client := <-s.Register:
			s.Clients[client.From] = client
			client.Status = true
		case client := <-s.UnRegister:
			if _, ok := s.Clients[client.From]; ok {
				delete(s.Clients, client.From)
				close(client.Send)
				client.Status = false
			}
		case message := <-s.Broadcast:
			if message.Type == "privet" {
				client := s.Clients[message.To]
				if client == nil{
					continue
				}
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.Clients, client.From)
				}
			} else if message.Type == "group" {
				var username []string
			services.CacheRun("users" , username)
			if username == nil {
				models.DB.Table("group_chats").Select("users.user_name").Where("group_chats.title = ?", message.To).
					Joins("left join user_groups on group_chats.id = user_groups.group_chat_id").
					Joins("left join users on user_groups.user_id = users.id").Find(&username)
				services.CacheRun("users" , username)
			}
					for _, e := range username {
						client := s.Clients[e]
						if client == nil{
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