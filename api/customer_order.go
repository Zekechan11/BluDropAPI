package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CustomerOrder struct {
	ID                int     `db:"Id" json:"id"`
	CustomerID        int     `db:"customer_id" json:"customer_id"`
	CustomerFullname  string  `db:"fullname" json:"customer_fullname"`
	AreaID            string  `db:"area_id" json:"area_id"`
	Num_gallons_order int     `db:"num_gallons_order" json:"num_gallons_order"`
	ReturnedGallons   int     `db:"returned_gallons" json:"returned_gallons"`
	COL 			  *int	  `db:"col" json:"col"`
	TotalCOL		  *int	  `db:"total_containers_on_loan" json:"total_containers_on_loan"`
	Date              string  `db:"date" json:"date"`
	Date_created      string  `db:"date_created" json:"date_created"`
	Total_price       float64 `db:"total_price" json:"total_price"`
	Payment           float64 `db:"payment" json:"payment"`
	PayableAmount	  *float64 `db:"payable_amount" json:"payable_amount"`
	Status            string  `db:"status" json:"status"`
	Agent			  *string  `db:"agent" json:"agent"`
}

func Customer_OrderRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_order", func(ctx *gin.Context) {
		area_id := ctx.DefaultQuery("area_id", "")
		status := ctx.DefaultQuery("status", "")

		var orders []CustomerOrder
		query := `
			SELECT 
				co.Id, 
				co.customer_id, 
				CONCAT(a.firstname, ' ', a.lastname) AS fullname,
				a.area_id,
				co.num_gallons_order,
				co.returned_gallons,
				(co.num_gallons_order - co.returned_gallons) AS col,
				lo.total_containers_on_loan,
				co.date, 
				co.date_created,
				co.total_price,
				co.payment,
				(co.total_price - co.payment) AS payable_amount,
				co.status,
				CONCAT(s.firstname, ' ', s.lastname) AS agent
			FROM 
				customer_order co
			LEFT JOIN 
				account_clients a ON co.customer_id = a.client_id
			LEFT JOIN
				account_staffs s ON co.area_id = s.area_id
			LEFT JOIN
				containers_on_loan lo ON co.customer_id = lo.customer_id
			WHERE
				s.role = 'Agent'
				AND (a.area_id = ? OR ? = '')
   				AND (co.status = ? OR ? = '')
		`
		err := db.Select(&orders, query, area_id, area_id, status, status)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, orders)
	})

	r.POST("/api/save_order", func(ctx *gin.Context) {
		// Log incoming request
		log.Println("Received save order request")

		// Struct to bind JSON input
		var insertCustomerOrder struct {
			CustomerID      int    `json:"customer_id"`
			NumGallonsOrder int    `json:"num_gallons_order"`
			Date            string `json:"date"`
			Status          string `json:"status"`
			AreaID          int    `json:"area_id"`
			Type            string `json:"type"`
		}

		// Bind JSON and log any binding errors
		if err := ctx.ShouldBindJSON(&insertCustomerOrder); err != nil {
			log.Printf("JSON Binding Error: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input",
				"details": err.Error(),
			})
			return
		}

		// Log received data for debugging
		log.Printf("Received Order Data: %+v", insertCustomerOrder)

		// Validate customer ID
		if insertCustomerOrder.CustomerID == 0 {
			log.Println("Error: Customer ID is NULL")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Customer ID is required",
			})
			return
		}

		// Validate date
		if insertCustomerOrder.Date == "" {
			insertCustomerOrder.Date = time.Now().Format("2006-01-02")
		}

		// Calculate total price based on inventory price
		var totalPrice float64
		getPriceQuery := fmt.Sprintf(`
			SELECT %s * ? 
			FROM pricing
			WHERE pricing_id = 1
		`, insertCustomerOrder.Type)
		err := db.Get(&totalPrice, getPriceQuery, insertCustomerOrder.NumGallonsOrder)
		if err != nil {
			log.Printf("Error calculating total price: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to calculate price",
				"details": err.Error(),
			})
			return
		}

		// Default status if not provided
		if insertCustomerOrder.Status == "" {
			insertCustomerOrder.Status = "Pending"
		}

		// Start a transaction
		tx, err := db.Beginx()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to start transaction",
				"details": err.Error(),
			})
			return
		}
		defer tx.Rollback() // Rollback in case of any error

		// Prepare insert query for customer order
		insertQuery := `
			INSERT INTO customer_order 
			(customer_id, num_gallons_order, date, date_created, total_price, status, area_id) 
			VALUES (?, ?, ?, NOW(), ?, ?, ?)`

		// Execute the query
		result, err := tx.Exec(insertQuery,
			insertCustomerOrder.CustomerID,
			insertCustomerOrder.NumGallonsOrder,
			insertCustomerOrder.Date,
			totalPrice,
			insertCustomerOrder.Status,
			insertCustomerOrder.AreaID,
		)

		if err != nil {
			log.Printf("Database Insertion Error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to save order",
				"details": err.Error(),
			})
			return
		}

		// Get last inserted ID
		lastID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last insert ID: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve order ID",
				"details": err.Error(),
			})
			return
		}

		// If status is Pending, subtract from inventory
		if insertCustomerOrder.Status == "Pending" {
			updateInventoryQuery := `
				UPDATE inventory_available 
				SET total_quantity = total_quantity - ?, 
				    last_updated = NOW()
				WHERE inventory_id = (
					SELECT inventory_id 
					FROM inventory_available 
					ORDER BY last_updated DESC 
					LIMIT 1
				)
			`
			_, err = tx.Exec(updateInventoryQuery, insertCustomerOrder.NumGallonsOrder)
			if err != nil {
				log.Printf("Error updating inventory: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to update inventory",
					"details": err.Error(),
				})
				return
			}
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to save order and update inventory",
				"details": err.Error(),
			})
			return
		}

		// Successful response
		log.Printf("Order saved successfully. ID: %d", lastID)
		ctx.JSON(http.StatusOK, gin.H{
			"message":     "Order saved successfully",
			"order_id":    lastID,
			"total_price": totalPrice,
			"order":       insertCustomerOrder,
		})
	})
}
