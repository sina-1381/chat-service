package chat

import (
	"chat-service/models"
	"chat-service/mongo"
	"log"
	"net/http"
	"time"

	"github.com/Kamva/mgm"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// websocket setup. it converts http to websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,                                       // read buffer
	WriteBufferSize:   1024,                                       // write buffer
	Subprotocols:      []string{"tcp"},                            // protocol tcp
	CheckOrigin:       func(r *http.Request) bool { return true }, // disable browser origin protection
	EnableCompression: true,                                       // compression enabled
}

// this is entrypoint of the websocket. it is registered as a route in main file
func Runner(se *Session, c *gin.Context) {
	// websocket connection stablished with the config set above
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// get user username from jwt token
	from := c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["username"]

	// creates a client for the user based on username and registers it in the websocket session channel
	client := &Client{Hub: se, Conn: conn, Send: make(chan *Message), From: from}
	client.Hub.Register <- client

	// registerd user client starts to read and write at the same time ready for
	// any message that is going to be send or recived by this user
	go client.Write()
	go client.Read()
}

// this is the read functionality of the client that handles the messages that user recives
func (c *Client) Read() {

	// if this function is stopped for any reason user in going to be unregisterd from websocket session
	defer func() {
		c.Hub.UnRegister <- c
	}()

	var message *Message

	for {
		err := c.Conn.ReadJSON(&message)

		// this handels unexpected errors. it eather logs them or breaks the function based on different errors
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message.From = c.From

		// open broadcast channel is open and ready to recive messages for the user client
		c.Hub.Broadcast <- message

		// if user broadcast channel recives a message from websocket session:
		// first checks if the message type is privet or group message
		if message.Type == "privet" {
			go func() {

				// store the message in mongodb messages collection
				messages := mongo.NewMessage(message.Msg, message.To, message.From.(string), message.Type, "send")
				err := mgm.Coll(messages).Create(messages)
				if err != nil {
					panic(err)
				}
			}()
		} else if message.Type == "group" {
			go func() {

				// if its a group message get the user ids of all users currently on this group
				// and add them to users array of the group message collection in mongodb
				var username []string
				models.DB.Table("group_chats").Select("users.user_name").Where("group_chats.title = ?", message.To).
					Joins("left join user_groups on group_chats.id = user_groups.group_chat_id").
					Joins("left join users on user_groups.user_id = users.id").Find(&username)
				GroupMessage := mongo.NewGroupMessage(message.Msg, message.To, "send", message.From.(string), message.Type, username)
				err := mgm.Coll(GroupMessage).Create(GroupMessage)
				if err != nil {
					panic(err)
				}
			}()
		} else {
			panic("wrong type !!!")
		}
	}
}

// this is the write functionality of the client that handles the messages that user sends to websocket session
func (c *Client) Write() {

	// on any error close this connection
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // 10 sec time limit to write message in conn
		err := c.Conn.WriteJSON(message)
		if err != nil {
			return
		}
	}

	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
