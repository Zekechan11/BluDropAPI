package dto

import "time"

type Message struct {
	MessageId int       `json:"message_id"`
	UserId    string    `json:"user_id"`
	AreaId    string    `json:"area_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
