package api

import (
	"log"
	"net/http"

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
		defer tx.Rollback() // Rollback in case of any error

		// Fetch the client order details using the new DTO
		var clientOrder dto.ClientOrder
		getOrderQuery := `SELECT total_price, num_gallons_order, area_id FROM customer_order WHERE Id = ?`
		err = tx.Get(&clientOrder, getOrderQuery, orderReq.CustomerID)
		if err != nil {
			log.Printf("Error fetching order: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch order details",
				"details": err.Error(),
			})
			return
		}

		// Update FGS count based on the fetched client order details
		updateFGSQuery := `UPDATE fgs SET count = count - ? WHERE area_id = ?`
		_, err = tx.Exec(updateFGSQuery, clientOrder.NumGallons, clientOrder.AreaID)
		if err != nil {
			log.Printf("Error updating FGS count: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update FGS count",
				"details": err.Error(),
			})
			return
		}

		// Update the order status and other details
		updateOrderQuery := `
			UPDATE customer_order 
			SET
				status = 'Completed',
				returned_gallons = ?,
				payment = ?
			WHERE id = ?
		`
		_, err = tx.Exec(updateOrderQuery, orderReq.GallonsToReturn, orderReq.Payment, orderReq.CustomerID)
		if err != nil {
			log.Printf("Error updating order status: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update order status",
				"details": err.Error(),
			})
			return
		}

		// Check if the customer exists in the containers_on_loan table
		var col dto.COL
		checkContainerQuery := `SELECT total_containers_on_loan, COUNT(*) FROM containers_on_loan WHERE customer_id = ?`
		err = tx.Get(&col, checkContainerQuery, orderReq.CustomerID)
		if err != nil {
			log.Printf("Error checking containers on loan: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check containers",
				"details": err.Error(),
			})
			return
		}

		// If no existing record, insert a new record into containers_on_loan
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
			// Update the existing record with the new number of containers on loan
			var previousNumGallons int
			if col.PreviousNumGallons != nil {
				previousNumGallons = *col.PreviousNumGallons
			} else {
				previousNumGallons = 0
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

		// Respond with success message
		ctx.JSON(http.StatusOK, gin.H{
			"message":    "Manual order processed successfully",
			"customerId": orderReq.CustomerID,
			"payment":    orderReq.Payment,
			"totalPrice": clientOrder.TotalPrice,
			"status":     "Completed",
		})
	})
}
