package util

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func ClientUsernameOrEmailCheck(db *sqlx.DB, username, email string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM account_clients
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

func SatffEmailCheck(db *sqlx.DB, email string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM account_staffs 
		WHERE email = ?
	`
	
	var count int
	err := db.Get(&count, query, email)
	if err != nil {
		log.Println("Error checking email existence:", err)
		return false, err
	}
	
	return count > 0, nil
}
