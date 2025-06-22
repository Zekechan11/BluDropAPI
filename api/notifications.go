package api

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Notification struct {
	DateCreated string  `db:"date_created" json:"date_created"`
	TotalPrice  float64 `db:"total_price" json:"total_price"`
	Payment     float64 `db:"payment" json:"payment"`
	Status      string  `db:"status" json:"status"`
	Unpaid      float64 `db:"unpaid" json:"unpaid"`
}

var (
	activeConnections int
	mu                sync.Mutex
	ticker            *time.Ticker
)

func NotificationRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/sse/notifications/:customer_id", func(c *gin.Context) {
		customerID := c.Param("customer_id")

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Flush()

		mu.Lock()
		activeConnections++
		if activeConnections == 1 {
			ticker = time.NewTicker(5 * time.Second)
			log.Println("First active connection detected, starting the ticker.")
		}
		mu.Unlock()

		defer func() {
			mu.Lock()
			activeConnections--
			if activeConnections == 0 {
				ticker.Stop()
				log.Println("No active connections left, stopping the ticker.")
			}
			mu.Unlock()
		}()

		for {
			select {
			case <-c.Writer.CloseNotify():
				log.Println("Client disconnected")
				return
			case <-ticker.C:
				var notifications []Notification
				query := `
					SELECT
						date_created,
						total_price,
						payment,
						status,
						SUM(total_price - payment) AS unpaid
					FROM customer_order
					WHERE customer_id = ?
						AND status = 'Pending'
						AND date_created <= DATE_ADD(CURDATE(), INTERVAL 1 MONTH)
					GROUP BY date_created, total_price, status, payment
					ORDER BY date_created DESC
				`

				err := db.Select(&notifications, query, customerID)
				if err != nil {
					log.Printf("DB error: %v", err)
					continue
				}

				jsonData, err := json.Marshal(notifications)
				if err != nil {
					log.Printf("Marshal error: %v", err)
					continue
				}

				fmt.Fprintf(c.Writer, "data: %s\n\n", jsonData)
				c.Writer.Flush()
			}
		}
	})
}
