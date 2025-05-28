package util

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func GenerateConversationID(userID int64) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%d%s", userID, timestamp)
}

func CreateChatID(db *sqlx.DB, userID int64, areaID int, staffOnly *bool) (string, error) {
	conversationID := GenerateConversationID(userID)

	var query string
	var args []interface{}

	if staffOnly != nil && *staffOnly {
		query = `
			INSERT INTO conversations (conversation_id, user_id, area_id, staff_only)
			VALUES (?, ?, ?, ?)
		`
		args = []interface{}{conversationID, userID, areaID, true}
	} else {
		query = `
			INSERT INTO conversations (conversation_id, user_id, area_id)
			VALUES (?, ?, ?)
		`
		args = []interface{}{conversationID, userID, areaID}
	}

	_, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Error creating conversation:", err)
		return "", err
	}

	return conversationID, nil
}
