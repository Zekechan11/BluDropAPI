package api

import (
	"bludrop-api/util"
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
	Type            string  `json:"type"`
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

		// prevent over pay
		overpay := 0.0
		paymentToInsert := orderReq.Payment
		if orderReq.Payment > totalPrice {
			overpay = orderReq.Payment - totalPrice
			paymentToInsert = totalPrice
		}

		if overpay > 0 {
			remainingOverpay, err := util.ApplyOverpay(tx, orderReq.CustomerID, overpay, nil)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to apply overpay to pending orders",
					"details": err.Error(),
				})
				return
			}
			overpay = remainingOverpay
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

		err = util.UpdateOrInsertContainersOnLoan(tx, orderReq.CustomerID, orderReq.GallonsToOrder, orderReq.GallonsToReturn)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to update containers on loan",
				"details": err.Error(),
			})
			return
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
