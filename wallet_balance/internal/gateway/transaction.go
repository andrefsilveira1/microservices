package gateway

import (
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
)

type TransactionGateway interface {
	Find(id string) (*entity.Transaction, error)
	Register(*entity.Transaction) error
}
