package seed

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {

	_, err := db.Exec(`Create table if not exists transactions_balance (id varchar(255),  account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date)`)
	if err != nil {
		log.Fatalf("Error creating clients table: %v", err)
	}
}
