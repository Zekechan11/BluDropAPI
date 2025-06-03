package api

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"bludrop-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	Conn           *websocket.Conn
	ConversationID string
}

var Clients = make(map[*websocket.Conn]*Client)
var ClientsMutex sync.RWMutex
var Broadcast = make(chan dto.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleMessages(db *sqlx.DB) {
	for {
		msg := <-Broadcast

		ClientsMutex.RLock()
		for conn, client := range Clients {
			if client.ConversationID == msg.ConversationId {
				err := conn.WriteJSON(msg)
				if err != nil {
					fmt.Println("Error writing message to client:", err)
					conn.Close()
					delete(Clients, conn)
				}
			}
		}
		ClientsMutex.RUnlock()
	}
}

func ChatRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/chat/:convo_id/messages", func(c *gin.Context) {
		convoId := c.Param("convo_id")

		query := `
			SELECT
				message_id,
				sender_id,
				conversation_id,
				sender_name,
				role,
				content,
				timestamp 
			FROM messages
			WHERE conversation_id = ?
			ORDER BY timestamp ASC`

		var messages []dto.MessageEntity
		err := db.Select(&messages, query, convoId)
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})

	r.GET("/chat/client/:id/messages", func(c *gin.Context) {
		clientId := c.Param("id")

		var convoId string
		err := db.Get(&convoId, "SELECT conversation_id FROM conversations WHERE user_id = ? AND staff_only = false", clientId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation ID: " + err.Error()})
			return
		}

		query := `
			SELECT
				message_id,
				sender_id,
				conversation_id,
				sender_name,
				role,
				content,
				timestamp 
			FROM messages
			WHERE conversation_id = ?
			ORDER BY timestamp ASC`

		var messages []dto.MessageEntity
		err = db.Select(&messages, query, convoId)
		if err != nil {
			fmt.Println("DB select error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"messages": messages, "convo_id": convoId})
	})

	r.GET("/chat/conversation/:area_id", func(c *gin.Context) {
		areaIdStr := c.Param("area_id")
		areaId, err := strconv.Atoi(areaIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area_id"})
			return
		}

		agentIdStr := c.Query("uid")
		agentId, err := strconv.Atoi(agentIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid uid"})
			return
		}

		query := `
			WITH admin_info AS (
				SELECT CONCAT(firstname, ' ', lastname) AS fullname, role
				FROM account_staffs
				WHERE role = 'Admin'
				LIMIT 1
			),

			latest_messages AS (
				SELECT m1.*
				FROM messages m1
				JOIN (
					SELECT conversation_id, MAX(timestamp) AS latest_timestamp
					FROM messages
					GROUP BY conversation_id
				) latest
				ON m1.conversation_id = latest.conversation_id
				AND m1.timestamp = latest.latest_timestamp
			)

			SELECT
				m.message_id,
				m.sender_id,
				c.conversation_id,
				m.sender_name AS last_sender,
				m.content,
				m.timestamp,
				ar.area AS area_name,
				CONCAT(a.firstname, ' ', a.lastname) AS fullname,
				a.role,
				a.type AS customer_type
			FROM conversations c
			LEFT JOIN account_clients a ON c.user_id = a.client_id
			LEFT JOIN areas ar ON a.area_id = ar.id
			LEFT JOIN latest_messages m ON c.conversation_id = m.conversation_id
			WHERE c.staff_only = false AND c.area_id = ?

			UNION ALL

			SELECT
				m.message_id,
				m.sender_id,
				c.conversation_id,
				m.sender_name AS last_sender,
				m.content,
				m.timestamp,
				NULL AS area_name,
				ai.fullname,
				ai.role,
				NULL AS customer_type
			FROM conversations c
			JOIN admin_info ai
			LEFT JOIN latest_messages m ON c.conversation_id = m.conversation_id
			WHERE c.staff_only = true AND c.user_id = ?

			ORDER BY timestamp DESC`

		var conversations []dto.ConversationEntity
		err = db.Select(&conversations, query, areaId, agentId)
		if err != nil {
			fmt.Println("DB query error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"conversations": conversations})
	})

	r.GET("/chat/conversation", func(c *gin.Context) {
		query := `
			WITH latest_messages AS (
				SELECT m1.*
				FROM messages m1
				JOIN (
					SELECT conversation_id, MAX(timestamp) AS latest_timestamp
					FROM messages
					GROUP BY conversation_id
				) latest ON m1.conversation_id = latest.conversation_id
						AND m1.timestamp = latest.latest_timestamp
			),
			users AS (
				SELECT
					ac.client_id AS user_id,
					CONCAT(ac.firstname, ' ', ac.lastname) AS fullname,
					ac.role,
					ac.type AS customer_type,
					ar.area AS area_name,
					'client' AS user_type
				FROM account_clients ac
				LEFT JOIN areas ar ON ac.area_id = ar.id

				UNION ALL

				SELECT
					ast.staff_id AS user_id,
					CONCAT(ast.firstname, ' ', ast.lastname) AS fullname,
					ast.role,
					NULL AS customer_type,
					ar.area AS area_name,
					'staff' AS user_type
				FROM account_staffs ast
				LEFT JOIN areas ar ON ast.area_id = ar.id
			)

			SELECT
				m.message_id,
				m.sender_id,
				c.conversation_id,
				m.sender_name AS last_sender,
				m.content,
				m.timestamp,
				u.area_name,
				u.fullname,
				u.role,
				u.customer_type
			FROM conversations c
			JOIN users u ON c.user_id = u.user_id AND (
				(u.user_type = 'client' AND c.staff_only = FALSE) OR
				(u.user_type = 'staff' AND c.staff_only = TRUE)
			)
			LEFT JOIN latest_messages m ON c.conversation_id = m.conversation_id
			ORDER BY m.timestamp DESC`

		var conversations []dto.ConversationEntity
		err := db.Select(&conversations, query)
		if err != nil {
			fmt.Println("DB query error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"conversations": conversations})
	})

	r.GET("/chat", func(c *gin.Context) {
		ConversationID := c.Query("convo_id")
		if ConversationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Conversation not Found"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}

		client := &Client{
			Conn:           conn,
			ConversationID: ConversationID,
		}

		ClientsMutex.Lock()
		Clients[conn] = client
		ClientsMutex.Unlock()

		defer func() {
			ClientsMutex.Lock()
			delete(Clients, conn)
			ClientsMutex.Unlock()
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
            INSERT INTO messages (sender_id, sender_name, role, content, conversation_id, timestamp)
            VALUES (:sender_id, :sender_name, :role, :content, :conversation_id, :timestamp)`

			_, err = db.NamedExec(query, map[string]any{
				"sender_id":       msg.SenderId,
				"sender_name":     msg.SenderName,
				"role":            msg.Role,
				"content":         msg.Content,
				"conversation_id": ConversationID,
				"timestamp":       time.Now(),
			})
			if err != nil {
				conn.WriteJSON(gin.H{"error": "Failed to store message"})
				continue
			}
			msg.ConversationId = ConversationID
			Broadcast <- msg
		}
	})
}
