package api

import (
	"bludrop-api/dto"
	"bludrop-api/util"
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

		// Fetch the client order using the dto.ClientOrder structure
		var clientOrder dto.ClientOrder
		getOrderQuery := `SELECT id, total_price, num_gallons_order, area_id FROM customer_order WHERE Id = ?`
		err = tx.Get(&clientOrder, getOrderQuery, paymentReq.OrderID)
		if err != nil {
			log.Printf("Error fetching order: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch order details",
				"details": err.Error(),
			})
			return
		}

		// Prevent overpaying
		overpay := 0.0
		paymentToInsert := paymentReq.AmountPaid
		if paymentReq.AmountPaid > clientOrder.TotalPrice{
			overpay = paymentReq.AmountPaid - clientOrder.TotalPrice
			paymentToInsert = clientOrder.TotalPrice
		}

		if overpay > 0 {
			remainingOverpay, err := util.ApplyOverpay(tx, paymentReq.CustomerID, overpay, &clientOrder.OrderID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to apply overpay to pending orders",
					"details": err.Error(),
				})
				return
			}
			overpay = remainingOverpay
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

		err = util.UpdateOrInsertContainersOnLoan(tx, paymentReq.CustomerID, clientOrder.NumGallons, paymentReq.GallonsReturned)
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
