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

func HandleMessages(db *sqlx.DB) {
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

func ChatRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/send_message", func(c *gin.Context) {
		var msg dto.Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
			return
		}

		// Insert message into the database using sqlx's NamedExec
		query := "INSERT INTO messages (user_id, area_id, content, timestamp) VALUES (:user_id, :area_id, :content, :timestamp)"
		_, err := db.NamedExec(query, map[string]interface{}{
			"user_id":    msg.SenderId,
			"area_id": msg.AreaId,
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

		query := `
			SELECT user_id, area_id, content, timestamp 
			FROM messages 
			WHERE user_id = :conversation_id OR area_id = :conversation_id 
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

	r.GET("/chat/customer/:username", func(ctx *gin.Context) {
		username := ctx.Param("username")

		query := `
			SELECT
				message_id,
				m.customer,
				m.sender_id,
				CONCAT(c.firstname, ' ', c.lastname) AS fullname,
				area_id,
				content,
				timestamp 
			FROM messages m
			LEFT JOIN client_accounts c ON m.sender_id = c.client_id
			WHERE m.customer = ?
			ORDER BY timestamp ASC`
		
		var messages []dto.MessageEntity
		err := db.Select(&messages, query, username)
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.GET("/chat/list/agent/:area_id", func(ctx *gin.Context) {
		areaId := ctx.Param("area_id")

		query := `
			SELECT
				message_id,
				m.customer,
				m.sender_id,
				CONCAT(c.firstname, ' ', c.lastname) AS fullname,
				area_id,
				content,
				timestamp 
			FROM messages m
			LEFT JOIN client_accounts c ON m.sender_id = c.client_id
			WHERE m.area_id = ?
			GROUP BY m.customer
			ORDER BY timestamp ASC`
		
		var messages []dto.MessageEntity
		err := db.Select(&messages, query, areaId)
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.GET("/chat", func(c *gin.Context) {
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

			query := `
				INSERT INTO messages (sender_id, area_id, customer, content, timestamp)
				VALUES (:sender_id, :area_id, :customer, :content, :timestamp)`

			_, err = db.NamedExec(query, map[string]interface{}{
				"sender_id":    msg.SenderId,
				"area_id": msg.AreaId,
				"customer": msg.Customer,
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
