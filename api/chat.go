package api

import (
	"fmt"
	"net/http"
	"time"
	"waterfalls/dto"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan dto.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/send_message", func(c *gin.Context) {
		var msg dto.Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
			return
		}

		// Insert message into the database using sqlx's NamedExec
		query := "INSERT INTO messages (sender, recipient, content, timestamp) VALUES (:sender, :recipient, :content, :timestamp)"
		_, err := db.NamedExec(query, map[string]interface{}{
			"sender":    msg.Sender,
			"recipient": msg.Recipient,
			"content":   msg.Content,
			"timestamp": time.Now(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store message"})
			return
		}

		// Broadcast the message to all connected clients
		Broadcast <- msg

		c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
	})

	r.GET("/get_message/:id", func(c *gin.Context) {
		conversationID := c.Param("id")

		// Use sqlx to query messages, and automatically map the result to a slice of dto.Message
		query := `
			SELECT sender, recipient, content, timestamp 
			FROM messages 
			WHERE sender = :conversation_id OR recipient = :conversation_id 
			ORDER BY timestamp ASC`
		
		var messages []dto.Message
		err := db.Select(&messages, query, map[string]interface{}{
			"conversation_id": conversationID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
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

			// Insert the message into the database
			query := "INSERT INTO messages (sender, recipient, content, timestamp) VALUES (:sender, :recipient, :content, :timestamp)"
			_, err = db.NamedExec(query, map[string]interface{}{
				"sender":    msg.Sender,
				"recipient": msg.Recipient,
				"content":   msg.Content,
				"timestamp": time.Now(),
			})
			if err != nil {
				conn.WriteJSON(gin.H{"error": "Failed to store message"})
				continue
			}

			// Broadcast the message
			Broadcast <- msg
		}
	})
}

func HandleMessages(db *sqlx.DB) {
	for {
		msg := <-Broadcast

		// Set the timestamp to the current time before broadcasting
		msg.Timestamp = time.Now()

		// Send the message to all connected clients
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
