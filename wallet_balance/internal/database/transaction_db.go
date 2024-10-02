package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
)

type TransactionDB struct {
	DB *sql.DB
}

func NewTransactionDB(db *sql.DB) *TransactionDB {
	return &TransactionDB{
		DB: db,
	}
}

func (t *TransactionDB) Register(transasction *entity.Transaction) error {
	fmt.Println("TRANSACTIONS RECEIVED ===>", transasction.ID)
	fmt.Println("TRANSACTIONS RECEIVED ===>", transasction.AccountIDFrom)
	fmt.Println("TRANSACTIONS RECEIVED ===>", transasction.AccountIDTo)
	fmt.Println("TRANSACTIONS RECEIVED ===>", transasction.Amount)
	stmt, err := t.DB.Prepare("INSERT INTO transactions_balance (id, account_id_from, account_id_to, amount, created_at) VALUES (?,?,?,?,?)")

	if err != nil {
		return err
	}
	fmt.Println("REGISTER ===>", err)
	defer stmt.Close()

	_, err = stmt.Exec(transasction.ID, transasction.AccountIDFrom, transasction.AccountIDTo, transasction.Amount, transasction.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDB) Find(id string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	fmt.Println("ID FIND ===>", id)
	stmt, err := t.DB.Prepare("SELECT id, account_id_from, account_id_to, amount, created_at from transactions_balance WHERE id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&transaction.ID,
		&transaction.AccountIDFrom,
		&transaction.AccountIDTo,
		&transaction.Amount,
		&transaction.CreatedAt,
	)

	fmt.Println("transaction transaction ===>", transaction.Amount)
	fmt.Println("transaction transaction ===>", transaction.AccountIDFrom)
	fmt.Println("transaction transaction ===>", transaction.AccountIDTo)

	if err != nil {
		return nil, err
	}

	log.Printf("Transaction found: ID=%s, AccountIDFrom=%s, AccountIDTo=%s, Amount=%f", transaction.ID, transaction.AccountIDFrom, transaction.AccountIDTo, transaction.Amount)
	return transaction, nil

}
