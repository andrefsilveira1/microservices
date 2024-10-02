package gateway

import "github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"

type BalancesGateway interface {
	UpdateBalances(id string, balance float64) error
	Find(id string) (*entity.Account, error)
}
