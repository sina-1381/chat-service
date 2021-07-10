package chat

import (
	"ginGorm/models"
	"ginGorm/mongo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kamva/mgm"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	Subprotocols:      []string{"tcp"},
	CheckOrigin: func(r *http.Request) bool {return true},
	EnableCompression: true,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func Runner(se *Session , c *gin.Context)  {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	from := c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["username"]
	client := &Client{Hub: se, Conn: conn, Send: make(chan *Message),From: from}
	client.Hub.Register <- client

	go client.Write()
	go client.Read()
}

func (c *Client)Read() {
	defer func() {
		c.Hub.UnRegister <- c
	}()
	var message *Message
	for {
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message.From = c.From
		c.Hub.Broadcast <- message
		if message.Type == "privet" {
			go func() {
				message := mongo.NewMessage(message.Msg, message.To, message.From.(string), message.Type, "send")
				err := mgm.Coll(message).Create(message)
				if err != nil {
					panic(err)
				}
			}()
		}else if message.Type == "group"{
			go func() {
				var username []string
				models.DB.Table("group_chats").Select("users.user_name").Where("group_chats.title = ?", message.To).
					Joins("left join user_groups on group_chats.id = user_groups.group_chat_id").
					Joins("left join users on user_groups.user_id = users.id").Find(&username)
				GroupMessage := mongo.NewGroupMessage(message.Msg , message.To , "send", message.From.(string) , message.Type ,username)
				err := mgm.Coll(GroupMessage).Create(GroupMessage)
				if err != nil {
					panic(err)
				}
			}()
		}else {
			panic("wrong type !!!")
		}
	}
}

func (c *Client)Write() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {

				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}
}