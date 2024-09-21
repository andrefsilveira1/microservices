package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            string
	AccountIDFrom string
	AccountIDTo   string
	Amount        float64
	CreatedAt     time.Time
}

func NewTransaction(account_id_from string, account_id_to string, amount float64) (*Transaction, error) {
	transaction := &Transaction{
		ID:            uuid.New().String(),
		AccountIDFrom: account_id_from,
		AccountIDTo:   account_id_to,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}

	return transaction, nil
	// Implement validate later
}
