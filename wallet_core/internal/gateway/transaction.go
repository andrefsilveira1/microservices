package gateway

import "github.com/andrefsilveira1/microservices/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
