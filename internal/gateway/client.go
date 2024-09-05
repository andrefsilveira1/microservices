package gateway

import "github.com/andrefsilveira1/microservices/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Add(client *entity.Client) error
}
