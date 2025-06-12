package api

import (
	"bludrop-api/dto"
	"database/sql"
	"fmt"
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
	GallonsToOrder  int     `json:"gallonsToOrder"`
	Type            string  `json:"type"`
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

		// Calculate total price
		var totalPrice float64
		getPriceQuery := fmt.Sprintf(`
			SELECT %s * ? 
			FROM pricing
			WHERE pricing_id = 1
		`, paymentReq.Type)
		err = tx.Get(&totalPrice, getPriceQuery, paymentReq.GallonsToOrder)
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
		paymentToInsert := paymentReq.AmountPaid
		if paymentReq.AmountPaid > totalPrice {
			overpay = paymentReq.AmountPaid - totalPrice
			paymentToInsert = totalPrice
		}

		// Fetch the client order using the dto.ClientOrder structure
		var clientOrder dto.ClientOrder
		getOrderQuery := `SELECT total_price, num_gallons_order, area_id FROM customer_order WHERE Id = ?`
		err = tx.Get(&clientOrder, getOrderQuery, paymentReq.OrderID)
		if err != nil {
			log.Printf("Error fetching order: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch order details",
				"details": err.Error(),
			})
			return
		}

		// Update FGS Count using the area_id from the client order
		updateFGSQuery := `UPDATE fgs SET count = count - ? WHERE area_id = ?`
		_, err = tx.Exec(updateFGSQuery, clientOrder.NumGallons, clientOrder.AreaID)
		if err != nil {
			log.Printf("Error updating fgs count: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update fgs count",
				"details": err.Error(),
			})
			return
		}

		// Update order status and payment details using the PaymentRequest
		updateOrderQuery := `
			UPDATE customer_order 
			SET
				status = 'Completed',
				returned_gallons = ?,
				payment = ?
			WHERE id = ?
		`
		_, err = tx.Exec(updateOrderQuery, paymentReq.GallonsReturned, paymentToInsert, paymentReq.OrderID)
		if err != nil {
			log.Printf("Error updating order status: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update order status",
				"details": err.Error(),
			})
			return
		}

		// Check if customer exists in containers_on_loan table using the COL structure
		var previousGallons int
		checkContainerQuery := `
			SELECT total_containers_on_loan 
			FROM containers_on_loan 
			WHERE customer_id = ?
			LIMIT 1
		`
		err = tx.Get(&previousGallons, checkContainerQuery, paymentReq.CustomerID)

		if err != nil {
			if err == sql.ErrNoRows {
				// No existing record, insert new
				insertContainerQuery := `
					INSERT INTO containers_on_loan 
					(customer_id, total_containers_on_loan, gallons_returned) 
					VALUES (?, ?, 0)
				`
				_, err = tx.Exec(insertContainerQuery, paymentReq.CustomerID, clientOrder.NumGallons)
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
			newNumGallons := previousGallons - paymentReq.GallonsReturned + clientOrder.NumGallons

			updateContainersQuery := `
				UPDATE containers_on_loan
				SET
					gallons_returned = ?,
					total_containers_on_loan = ?
				WHERE customer_id = ?
			`
			_, err = tx.Exec(updateContainersQuery, paymentReq.GallonsReturned, newNumGallons, paymentReq.CustomerID)
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
			"totalPrice": clientOrder.TotalPrice,
			"overpay":    overpay,
			"status":     "Completed",
		})
	})
}
