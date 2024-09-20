package entity

import (
	"time"
)

type Transaction struct {
	ID          string
	AccountFrom *string
	AccountTo   *string
	Amount      float64
	CreatedAt   time.Time
}
