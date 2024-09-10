package createaccount

import (
	"github.com/andrefsilveira1/microservices/internal/entity"
	"github.com/andrefsilveira1/microservices/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientId string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(a gateway.AccountGateway, c gateway.ClientGateway) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: a,
		ClientGateway:  c,
	}
}

func (u *CreateAccountUseCase) Execute(input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	client, err := u.ClientGateway.Get(input.ClientId)
	if err != nil {
		return nil, err
	}

	account := entity.NewAccount(client)
	err = u.AccountGateway.Save(account)

	if err != nil {
		return nil, err
	}

	return &CreateAccountOutputDTO{
		ID: account.ID,
	}, nil
}
