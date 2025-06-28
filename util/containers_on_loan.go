package util

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func UpdateOrInsertContainersOnLoan(tx *sqlx.Tx, customerID int, gallonsToOrder int, gallonsReturned int) error {
	var previousGallons int

	// Check if the customer already has a record in the containers_on_loan table
	checkContainerQuery := `
		SELECT total_containers_on_loan 
		FROM containers_on_loan 
		WHERE customer_id = ? 
		LIMIT 1
	`
	err := tx.QueryRow(checkContainerQuery, customerID).Scan(&previousGallons)

	if err != nil {
		if err == sql.ErrNoRows {
			// No existing record, insert new one
			insertContainerQuery := `
				INSERT INTO containers_on_loan 
				(customer_id, total_containers_on_loan, gallons_returned) 
				VALUES (?, ?, 0)
			`
			_, err = tx.Exec(insertContainerQuery, customerID, gallonsToOrder)
			if err != nil {
				log.Printf("Error inserting containers_on_loan for customer %d: %v", customerID, err)
				return fmt.Errorf("failed to record containers on loan for customer %d: %v", customerID, err)
			}
		} else {
			log.Printf("Error checking containers on loan for customer %d: %v", customerID, err)
			return fmt.Errorf("failed to check containers on loan for customer %d: %v", customerID, err)
		}
	} else {
		// If record exists, update the existing record
		newNumGallons := previousGallons - gallonsReturned + gallonsToOrder

		updateContainersQuery := `
			UPDATE containers_on_loan
			SET
				gallons_returned = ?,
				total_containers_on_loan = ?
			WHERE customer_id = ?
		`
		_, err = tx.Exec(updateContainersQuery, gallonsReturned, newNumGallons, customerID)
		if err != nil {
			log.Printf("Error updating containers on loan for customer %d: %v", customerID, err)
			return fmt.Errorf("failed to update containers on loan for customer %d: %v", customerID, err)
		}
	}

	return nil
}
