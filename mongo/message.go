package mongo

import "github.com/Kamva/mgm"

// messages model in mongodb collection
type Message struct {
	mgm.DefaultModel `bson:",inline"`
	Msg              string `json:"msg" bson:"msg"`
	To               string `json:"to" bson:"to"`
	From             string `json:"from" bson:"from"`
	Type             string `json:"type" bson:"type"`
	Status           string `json:"status" bson:"status"`
}

// this func recives the data and creates a new message ready to be stored in mongodb
func NewMessage(message, to, from, tpe, status string) *Message {
	return &Message{
		Msg:    message,
		To:     to,
		From:   from,
		Type:   tpe,
		Status: status,
	}
}

// group messages model in mongodb collection
type GroupMessage struct {
	mgm.DefaultModel `bson:",inline"`
	Msg              string   `json:"msg" bson:"msg"`
	Title            string   `json:"title" bson:"title"`
	To               []string `json:"to" bson:"to"`
	From             string   `json:"from" bson:"from"`
	Type             string   `json:"type" bson:"type"`
	Status           string   `json:"status" bson:"status"`
}

// this func recives the data and creates a new group message ready to be stored in mongodb
func NewGroupMessage(message, title, status, from, tpe string, to []string) *GroupMessage {
	return &GroupMessage{
		Msg:    message,
		Title:  title,
		To:     to,
		From:   from,
		Type:   tpe,
		Status: status,
	}
}
