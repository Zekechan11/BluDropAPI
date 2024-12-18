package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type PaymentRequest struct {
	OrderID         int     `json:"orderId"`
	CustomerID      int     `json:"customerId"`
	AmountPaid      float64 `json:"amountPaid"`
	GallonsReturned int     `json:"gallonsReturned"`
}

func PaymentRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/api/process-payment", func(ctx *gin.Context) {
		var paymentReq PaymentRequest

		// Bind JSON input
		if err := ctx.ShouldBindJSON(&paymentReq); err != nil {
			log.Printf("JSON Binding Error: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input",
				"details": err.Error(),
			})
			return
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

		// Fetch the order to get total price
		var totalPrice float64
		getOrderQuery := `SELECT total_price FROM customer_order WHERE Id = ?`
		err = tx.Get(&totalPrice, getOrderQuery, paymentReq.OrderID)
		if err != nil {
			log.Printf("Error fetching order total price: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch order details",
				"details": err.Error(),
			})
			return
		}

		// Check if payment amount matches total price
		if paymentReq.AmountPaid < totalPrice {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient payment amount",
				"details": gin.H{
					"totalPrice": totalPrice,
					"amountPaid": paymentReq.AmountPaid,
					"difference": totalPrice - paymentReq.AmountPaid,
				},
			})
			return
		}

		// Update order status
		updateOrderQuery := `
			UPDATE customer_order 
			SET status = 'Completed' 
			WHERE Id = ?
		`
		_, err = tx.Exec(updateOrderQuery, paymentReq.OrderID)
		if err != nil {
			log.Printf("Error updating order status: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update order status",
				"details": err.Error(),
			})
			return
		}

		// Check if customer exists in containers_on_loan
		var existingRecord int
		checkContainerQuery := `
			SELECT COUNT(*) FROM containers_on_loan 
			WHERE customer_id = ?
		`
		err = tx.Get(&existingRecord, checkContainerQuery, paymentReq.CustomerID)
		if err != nil {
			log.Printf("Error checking existing containers: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check containers",
				"details": err.Error(),
			})
			return
		}

		// If no existing record, insert new record
		if existingRecord == 0 {
			// Get the number of gallons ordered
			var numGallons int
			getGallonsQuery := `
				SELECT num_gallons_order FROM customer_order 
				WHERE Id = ?
			`
			err = tx.Get(&numGallons, getGallonsQuery, paymentReq.OrderID)
			if err != nil {
				log.Printf("Error fetching order gallons: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to fetch order gallons",
					"details": err.Error(),
				})
				return
			}

			// Insert new record in containers_on_loan
			insertContainerQuery := `
				INSERT INTO containers_on_loan 
				(customer_id, total_containers_on_loan, gallons_returned) 
				VALUES (?, ?, 0)
			`
			_, err = tx.Exec(insertContainerQuery, paymentReq.CustomerID, numGallons)
			if err != nil {
				log.Printf("Error inserting containers_on_loan: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to record containers on loan",
					"details": err.Error(),
				})
				return
			}
		} else {
			// Update existing record with returned gallons
			updateContainersQuery := `
				UPDATE containers_on_loan 
				SET gallons_returned = gallons_returned + ? 
				WHERE customer_id = ?
			`
			_, err = tx.Exec(updateContainersQuery, paymentReq.GallonsReturned, paymentReq.CustomerID)
			if err != nil {
				log.Printf("Error updating gallons returned: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to update gallons returned",
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
				"error":   "Failed to process payment",
				"details": err.Error(),
			})
			return
		}

		// Successful response
		ctx.JSON(http.StatusOK, gin.H{
			"message":    "Payment processed successfully",
			"orderId":    paymentReq.OrderID,
			"amountPaid": paymentReq.AmountPaid,
			"totalPrice": totalPrice,
			"status":     "Completed",
		})
	})
}
