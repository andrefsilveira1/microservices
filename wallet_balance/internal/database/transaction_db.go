package database

import (
	"database/sql"

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

func (t *TransactionDB) Find(id string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	stmt, err := t.DB.Prepare("SELECT id, account_id_from, account_id_to, amount, created_at from transactions WHERE id = ?")
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
	if err != nil {
		return nil, err
	}

	return transaction, nil

}
