package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan dto.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/send_message", func(c *gin.Context) {
		var msg dto.Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
			return
		}

		query := "INSERT INTO messages (sender, recipient, content, timestamp) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(query, msg.Sender, msg.Recipient, msg.Content, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store message"})
			return
		}

		Broadcast <- msg

		c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
	})

	r.GET("/get_message/:id", func(c *gin.Context) {
		conversationID := c.Param("id")
	
		query := `
			SELECT sender, recipient, content, timestamp 
			FROM messages 
			WHERE sender = ? OR recipient = ? 
			ORDER BY timestamp ASC`
		
		rows, err := db.Query(query, conversationID, conversationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}
		defer rows.Close()
	
		var messages []dto.Message
		for rows.Next() {
			var sender, recipient, content, timestamp string
			if err := rows.Scan(&sender, &recipient, &content, &timestamp); err != nil {
				fmt.Println("Row scan error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan message data"})
				return
			}
		
			ts, err := time.Parse("2006-01-02 15:04:05", timestamp)
			if err != nil {
				fmt.Println("Timestamp parsing error:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse timestamp"})
				return
			}
		
			messages = append(messages, dto.Message{
				Sender:    sender,
				Recipient: recipient,
				Content:   content,
				Timestamp: ts,
			})
		}
		
	
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate through results"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})
	

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}

		Clients[conn] = true
		defer func() {
			delete(Clients, conn)
			conn.Close()
		}()

		for {
			var msg dto.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					break
				}
				fmt.Println("WebSocket read error:", err)
				break
			}

			query := "INSERT INTO messages (sender, recipient, content, timestamp) VALUES (?, ?, ?, ?)"
			_, err = db.Exec(query, msg.Sender, msg.Recipient, msg.Content, time.Now())
			if err != nil {
				conn.WriteJSON(gin.H{"error": "Failed to store message"})
				continue
			}

			Broadcast <- msg
		}
	})
}

func HandleMessages(db *sql.DB) {
	for {
		msg := <-Broadcast

		msg.Timestamp = time.Now()

		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("Error writing message to client:", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
