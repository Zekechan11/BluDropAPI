package api

import (
	"log"
	"net/http"
	"time"

	"waterfalls/dto" // Import the dto package

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ManualOrderRequest struct {
	CustomerID      int     `json:"customerId"`
	GallonsToOrder  int     `json:"gallonsToOrder"`
	Payment         float64 `json:"payment"`
	GallonsToReturn int     `json:"gallonsToReturn"`
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
		getPriceQuery := `
			SELECT price * ? 
			FROM inventory_available 
			ORDER BY last_updated DESC 
			LIMIT 1
		`
		err = tx.Get(&totalPrice, getPriceQuery, orderReq.GallonsToOrder)
		if err != nil {
			log.Printf("Error calculating total price: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to calculate price",
				"details": err.Error(),
			})
			return
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
			orderReq.Payment,
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
		var col dto.COL
		checkContainerQuery := `
			SELECT total_containers_on_loan, COUNT(*) 
			FROM containers_on_loan 
			WHERE customer_id = ?
		`
		err = tx.Get(&col, checkContainerQuery, orderReq.CustomerID)
		if err != nil {
			log.Printf("Error checking containers on loan: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check containers",
				"details": err.Error(),
			})
			return
		}

		if col.ExistingRecord == 0 {
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
			var previousNumGallons int
			if col.PreviousNumGallons != nil {
				previousNumGallons = *col.PreviousNumGallons
			}

			newNumGallons := previousNumGallons - orderReq.GallonsToReturn + orderReq.GallonsToOrder

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
			"payment":    orderReq.Payment,
			"totalPrice": totalPrice,
			"status":     "Completed",
		})
	})
}
