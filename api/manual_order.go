package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ManualOrderRequest struct {
	CustomerID      int     `json:"customerId"`
	GallonsToOrder  int     `json:"gallonsToOrder"`
	Payment         float64 `json:"payment"`
	GallonsToReturn int     `json:"gallonsToReturn"`
	Type            string `json:"type"`
}

func ManualOrderRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/api/process-manual-order", func(ctx *gin.Context) {
		var orderReq ManualOrderRequest

		// Bind JSON input
		if err := ctx.ShouldBindJSON(&orderReq); err != nil {
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
		defer tx.Rollback()

		// Calculate total price
		var totalPrice float64
		getPriceQuery := fmt.Sprintf(`
			SELECT %s * ? 
			FROM pricing
			WHERE pricing_id = 1
		`, orderReq.Type)
		err = tx.Get(&totalPrice, getPriceQuery, orderReq.GallonsToOrder)
		if err != nil {
			log.Printf("Error calculating total price: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to calculate price",
				"details": err.Error(),
			})
			return
		}

		// Prevent overpaying
		overpay := 0.0
		paymentToInsert := orderReq.Payment
		if orderReq.Payment > totalPrice {
			overpay = orderReq.Payment - totalPrice
			paymentToInsert = totalPrice
		}

		// Insert into customer_order
		insertOrderQuery := `
		INSERT INTO customer_order 
		(customer_id, num_gallons_order, date, date_created, total_price, payment, returned_gallons, status, area_id) 
		VALUES (?, ?, ?, NOW(), ?, ?, ?, 'Completed', 
			(SELECT area_id FROM account_clients WHERE client_id = ?))
		`
		_, err = tx.Exec(
			insertOrderQuery,
			orderReq.CustomerID,
			orderReq.GallonsToOrder,
			time.Now().Format("2006-01-02"),
			totalPrice,
			paymentToInsert,
			orderReq.GallonsToReturn,
			orderReq.CustomerID,
		)
		if err != nil {
			log.Printf("Error inserting order: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to insert order",
				"details": err.Error(),
			})
			return
		}

		// Update FGS Count
		updateFGSQuery := `
			UPDATE fgs 
			SET count = count - ? 
			WHERE area_id = (
				SELECT area_id 
				FROM account_clients 
				WHERE client_id = ?
			)
		`
		_, err = tx.Exec(updateFGSQuery, orderReq.GallonsToOrder, orderReq.CustomerID)
		if err != nil {
			log.Printf("Error updating fgs count: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update fgs count",
				"details": err.Error(),
			})
			return
		}

		// Handle containers_on_loan similarly to previous implementation
		var previousGallons int
		checkContainerQuery := `
			SELECT total_containers_on_loan 
			FROM containers_on_loan 
			WHERE customer_id = ?
			LIMIT 1
		`
		err = tx.Get(&previousGallons, checkContainerQuery, orderReq.CustomerID)

		if err != nil {
			if err == sql.ErrNoRows {
				// No existing record, insert new
				insertContainerQuery := `
					INSERT INTO containers_on_loan 
					(customer_id, total_containers_on_loan, gallons_returned) 
					VALUES (?, ?, 0)
				`
				_, err = tx.Exec(insertContainerQuery, orderReq.CustomerID, orderReq.GallonsToOrder)
				if err != nil {
					log.Printf("Error inserting containers_on_loan: %v", err)
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to record containers on loan",
						"details": err.Error(),
					})
					return
				}
			} else {
				log.Printf("Error checking containers on loan: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to check containers",
					"details": err.Error(),
				})
				return
			}
		} else {
			// Update the existing record
			newNumGallons := previousGallons - orderReq.GallonsToReturn + orderReq.GallonsToOrder

			updateContainersQuery := `
				UPDATE containers_on_loan
				SET
					gallons_returned = ?,
					total_containers_on_loan = ?
				WHERE customer_id = ?
			`
			_, err = tx.Exec(updateContainersQuery, orderReq.GallonsToReturn, newNumGallons, orderReq.CustomerID)
			if err != nil {
				log.Printf("Error updating containers on loan: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to update containers on loan",
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
				"error":   "Failed to process order",
				"details": err.Error(),
			})
			return
		}

		// Successful response
		ctx.JSON(http.StatusOK, gin.H{
			"message":    "Manual order processed successfully",
			"customerId": orderReq.CustomerID,
			"payment":    paymentToInsert,
			"totalPrice": totalPrice,
			"overpay":    overpay,
			"status":     "Completed",
		})
	})
}
