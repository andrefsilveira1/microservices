package seed

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS clients (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		created_at DATE
	)`)
	if err != nil {
		log.Fatalf("Error creating clients table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id VARCHAR(255) PRIMARY KEY,
		client_id VARCHAR(255),
		balance DECIMAL(10, 2),
		created_at DATE
	)`)
	if err != nil {
		log.Fatalf("Error creating accounts table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		id VARCHAR(255) PRIMARY KEY,
		account_id_from VARCHAR(255),
		account_id_to VARCHAR(255),
		amount INT,
		created_at DATE
	)`)
	if err != nil {
		log.Fatalf("Error creating transactions table: %v", err)
	}

	InsertFakeData(db)
}
