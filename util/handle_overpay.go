package util

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type PendingOrder struct {
	OrderID    int     `db:"id"`
	TotalPrice float64 `db:"total_price"`
	Payment    float64 `db:"payment"`
}

func ApplyOverpay(tx *sqlx.Tx, customerID int, overpay float64) (float64, error) {
	if overpay <= 0 {
		return 0, nil
	}

	var pendingOrders []PendingOrder
	getPendingOrdersQuery := `
		SELECT id, total_price, payment
		FROM customer_order
		WHERE customer_id = ? AND status = 'Pending'
		ORDER BY date_created ASC
	`

	err := tx.Select(&pendingOrders, getPendingOrdersQuery, customerID)
	if err != nil {
		log.Printf("Error fetching pending orders: %v", err)
		return overpay, err
	}

	for _, pending := range pendingOrders {
		if overpay <= 0 {
			break
		}

		remaining := pending.TotalPrice - pending.Payment
		if remaining <= 0 {
			continue
		}

		if overpay >= remaining {
			_, err := tx.Exec(`
				UPDATE customer_order
				SET payment = total_price, status = 'Completed'
				WHERE id = ?
			`, pending.OrderID)
			if err != nil {
				return overpay, err
			}
			overpay -= remaining
		} else {
			_, err := tx.Exec(`
				UPDATE customer_order
				SET payment = payment + ?
				WHERE id = ?
			`, overpay, pending.OrderID)
			if err != nil {
				return overpay, err
			}
			overpay = 0
		}
	}

	return overpay, nil
}
