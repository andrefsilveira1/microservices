package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            string    `json:"ID"`
	AccountIDFrom string    `json:"account_id_from"`
	AccountIDTo   string    `json:"account_id_to"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
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
}
