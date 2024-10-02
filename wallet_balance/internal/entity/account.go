package entity

import (
	"time"
)

type Account struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
