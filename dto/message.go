package dto

import "time"

type Conversation struct {
	ConversationId int `json:"conversation_id"`
	UserId         int `json:"user_id"`
	AreaId         int `json:"area_id"`
	StaffOnly      int `json:"staff_only"`
}

type Message struct {
	MessageId      int       `json:"message_id" db:"message_id"`
	SenderId       int       `json:"sender_id" db:"sender_id"`
	SenderName     string    `json:"sender_name" db:"sender_name"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	ConversationId string    `json:"conversation_id" db:"conversation_id"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
}

type MessageEntity struct {
	MessageId      int    `json:"message_id" db:"message_id"`
	SenderId       int    `json:"sender_id" db:"sender_id"`
	SenderName     string `json:"sender_name" db:"sender_name"`
	Role           string `json:"role" db:"role"`
	Content        string `json:"content" db:"content"`
	ConversationId int    `json:"conversation_id" db:"conversation_id"`
	Timestamp      string `json:"timestamp" db:"timestamp"`
}

type ConversationEntity struct {
	MessageId      *int    `json:"message_id" db:"message_id"`
	SenderId       *int    `json:"sender_id" db:"sender_id"`
	LastSender     *string `json:"last_sender" db:"last_sender"`
	Fullname       string  `json:"fullname" db:"fullname"`
	Content        *string `json:"content" db:"content"`
	AreaName       *string `json:"area_name" db:"area_name"`
	Role           *string `json:"role" db:"role"`
	CustomerType   *string `json:"customer_type" db:"customer_type"`
	ConversationId string  `json:"conversation_id" db:"conversation_id"`
	Timestamp      *string `json:"timestamp" db:"timestamp"`
}

type StartChat struct {
	UserId int64 `json:"user_id"`
	AreaId int   `json:"area_id"`
}
