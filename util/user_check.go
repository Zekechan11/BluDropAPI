package util

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func CheckUsernameOrEmailExists(db *sqlx.DB, username, email string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM client_accounts 
		WHERE username = ? OR email = ?
	`
	
	var count int
	err := db.Get(&count, query, username, email)
	if err != nil {
		log.Println("Error checking username/email existence:", err)
		return false, err
	}
	
	return count > 0, nil
}
