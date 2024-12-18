package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type CustomerOrder struct {
	ID                int    `db:"Id"`
	CustomerID        int    `db:"customer_id"`
	CustomerFirstName string `db:"FirstName"`
	CustomerLastName  string `db:"LastName"`
	Num_gallons_order int    `db:"num_gallons_order"`
	Date              string `db:"date"`
	Date_created      string `db:"date_created"`
}

func Customer_OrderRoutes(r *gin.Engine, db *sqlx.DB) {
	r.GET("/api/get_order", func(ctx *gin.Context) {
		var orders []CustomerOrder
		query := `
			SELECT 
				co.Id, 
				co.customer_id, 
				a.FirstName, 
				a.LastName, 
				co.num_gallons_order, 
				co.date, 
				co.date_created 
			FROM 
				customer_order co
			LEFT JOIN 
				Accounts a ON co.customer_id = a.Id
		`
		err := db.Select(&orders, query)
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
			CustomerID        string `json:"customer_id"`
			Num_gallons_order string `json:"num_gallons_order"`
			Date              string `json:"date"`
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
		if insertCustomerOrder.CustomerID == "" {
			log.Println("Error: Customer ID is NULL")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Customer ID is required",
			})
			return
		}

		// Convert num_gallons_order to int
		numGallons, err := strconv.Atoi(insertCustomerOrder.Num_gallons_order)
		if err != nil {
			log.Printf("Error converting num_gallons_order: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid number of gallons",
				"details": err.Error(),
			})
			return
		}

		// Validate date
		if insertCustomerOrder.Date == "" {
			insertCustomerOrder.Date = time.Now().Format("2006-01-02")
		}

		// Prepare insert query
		insertQuery := `
			INSERT INTO customer_order 
			(customer_id, num_gallons_order, date, date_created) 
			VALUES (?, ?, ?, NOW())`

		// Execute the query
		result, err := db.Exec(insertQuery,
			insertCustomerOrder.CustomerID,
			numGallons,
			insertCustomerOrder.Date,
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

		// Successful response
		log.Printf("Order saved successfully. ID: %d", lastID)
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Order saved successfully",
			"order_id": lastID,
			"order":    insertCustomerOrder,
		})
	})
}
