package database

import (
	"database/sql"
	"fmt"
	"time"

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
	stmt, err := t.DB.Prepare("INSERT INTO transactions_balance (id, account_id_from, account_id_to, amount, created_at) VALUES (?,?,?,?,?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(transasction.ID, transasction.AccountIDFrom, transasction.AccountIDTo, transasction.Amount, transasction.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDB) Find(id string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	stmt, err := t.DB.Prepare("SELECT id, account_id_from, account_id_to, amount, created_at from transactions_balance WHERE id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(id)
	var createdAtRaw []byte
	err = row.Scan(
		&transaction.ID,
		&transaction.AccountIDFrom,
		&transaction.AccountIDTo,
		&transaction.Amount,
		&createdAtRaw,
	)

	if err != nil {
		fmt.Println("Error during scan:", err)
		return nil, err
	}

	createdAtStr := string(createdAtRaw)
	transaction.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		fmt.Println("Error parsing created_at:", err)
		return nil, err
	}

	fmt.Println("Transaction found:", transaction)
	return transaction, nil

}
