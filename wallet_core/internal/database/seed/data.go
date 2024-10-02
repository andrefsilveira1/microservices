package seed

import (
	"database/sql"
	"log"
)

func InsertFakeData(db *sql.DB) {
	_, err := db.Exec(`INSERT INTO clients (id, name, email, created_at) VALUES 
		('1', 'John Doe', 'john@example.com', CURDATE()),
		('2', 'Jane Doe', 'jane@example.com', CURDATE())`)
	if err != nil {
		log.Printf("Error inserting fake data into clients: %v", err)
	}

	_, err = db.Exec(`INSERT INTO accounts (id, client_id, balance, created_at) VALUES 
		('1', '1', 1000, CURDATE()),
		('2', '2', 1500, CURDATE())`)
	if err != nil {
		log.Printf("Error inserting fake data into accounts: %v", err)
	}

	_, err = db.Exec(`INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES 
		('1', '1', '2', 200, CURDATE())`)
	if err != nil {
		log.Printf("Error inserting fake data into transactions: %v", err)
	}
}

func DropTables(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS transactions")
	if err != nil {
		log.Printf("Error dropping transactions table: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS accounts")
	if err != nil {
		log.Printf("Error dropping accounts table: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS clients")
	if err != nil {
		log.Printf("Error dropping clients table: %v", err)
	}

	log.Println("All tables dropped successfully.")
}
