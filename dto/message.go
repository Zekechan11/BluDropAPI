package dto

import "time"

type MessageEntity struct {
	MessageId int       `db:"message_id"`
	UserId    string    `db:"user_id"`
	FirstName    string    `db:"firstname"`
	LastName    string    `db:"lastname"`
	AreaId    string    `db:"area_id"`
	Content   string    `db:"content"`
	Timestamp string `db:"timestamp"`
}

type Message struct {
	MessageId int       `json:"message_id"`
	UserId    string    `json:"user_id"`
	AreaId    string    `json:"area_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
