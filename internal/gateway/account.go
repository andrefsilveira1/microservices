package gateway

import "github.com/andrefsilveira1/microservices/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	Find(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
