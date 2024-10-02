package seed

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions_balance (
			id VARCHAR(255),  
			account_id_from VARCHAR(255), 
			account_id_to VARCHAR(255), 
			amount DECIMAL(10, 2), 
			created_at DATETIME
		)
	`)
	if err != nil {
		log.Fatalf("Error creating transactions_balance table: %v", err)
	}
	log.Println("Table 'transactions_balance' created or already exists.")
}

func DropTables(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS transactions_balance`)
	if err != nil {
		log.Fatalf("Error dropping transactions_balance table: %v", err)
	}
	log.Println("Table 'transactions_balance' dropped.")
}
