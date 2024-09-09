package createtransaction

import (
	"github.com/andrefsilveira1/microservices/internal/entity"
	"github.com/andrefsilveira1/microservices/internal/gateway"
	"github.com/andrefsilveira1/microservices/pkg/events"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string
	AccountIDTo   string
	Amount        float64
}

type CreateTransactionOutputDTO struct {
	ID string
}

type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(transactionGateway gateway.TransactionGateway, accountGateway gateway.AccountGateway, eventDispatcher events.EventDispatcherInterface, transactionCreated events.EventInterface) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}

func (u *CreateTransactionUseCase) Execute(input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	accountFrom, err := u.AccountGateway.Find(input.AccountIDFrom)
	if err != nil {
		return nil, err
	}

	accountTo, err := u.AccountGateway.Find(input.AccountIDTo)
	if err != nil {
		return nil, err
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		return nil, err
	}

	err = u.TransactionGateway.Create(transaction)
	if err != nil {
		return nil, err
	}
	output := &CreateTransactionOutputDTO{ID: transaction.ID}
	u.TransactionCreated.GetPayload()
	u.EventDispatcher.Dispatch(u.TransactionCreated)

	return output, nil
}
