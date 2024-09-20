package entity

import (
	"time"
)

type Transaction struct {
	ID            string
	AccountIDFrom *string
	AccountIDTo   *string
	Amount        float64
	CreatedAt     time.Time
}
