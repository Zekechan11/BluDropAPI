package dto

import "time"

type MessageEntity struct {
	MessageId int    `json:"message_id" db:"message_id"`
	SenderId  string `json:"sender_id" db:"sender_id"`
	Customer  string `json:"customer" db:"customer"`
	Fullname  string `json:"fullname" db:"fullname"`
	AreaId    string `json:"area_id" db:"area_id"`
	Content   string `json:"content" db:"content"`
	Timestamp string `json:"timestamp" db:"timestamp"`
}

type Message struct {
	MessageId int       `json:"message_id"`
	SenderId  string    `json:"sender_id"`
	AreaId    string    `json:"area_id"`
	Customer  string    `json:"customer"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
