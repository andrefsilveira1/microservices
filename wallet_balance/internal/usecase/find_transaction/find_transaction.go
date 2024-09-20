package gettransaction

import "time"

type FindTransactionInputDTO struct {
	ID string `json:"id"`
}

type FindTransactionOutputDTO struct {
	ID          string    `json:"id"`
	AccountFrom *string   `json:"account_id_from"`
	AccountTo   *string   `json:"account_id_to"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}
