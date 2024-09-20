package gettransaction

import (
	"context"
	"time"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/gateway"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/uow"
)

type FindTransactionInputDTO struct {
	ID string `json:"id"`
}

type FindTransactionOutputDTO struct {
	ID            string    `json:"id"`
	AccountIDFrom *string   `json:"account_id_from"`
	AccountIDTo   *string   `json:"account_id_to"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type FindTransactionUseCase struct {
	Uow              uow.UowInterface
	EventDispatcher  events.EventDispatcherInterface
	TransactionFound events.EventInterface
}

func NewFindTransactionUseCase(uow uow.UowInterface, eventDispatcher events.EventDispatcherInterface, transactionFound events.EventInterface) *FindTransactionUseCase {
	return &FindTransactionUseCase{
		Uow:              uow,
		EventDispatcher:  eventDispatcher,
		TransactionFound: transactionFound,
	}
}

func (u *FindTransactionUseCase) Execute(ctx context.Context, input FindTransactionInputDTO) (*FindTransactionOutputDTO, error) {
	output := &FindTransactionOutputDTO{}
	err := u.Uow.Do(ctx, func(_ *uow.Uow) error {
		transactionRepository := u.getTransactionRepository(ctx)

		transaction, err := transactionRepository.Find(input.ID)
		if err != nil {
			return err
		}
		output.AccountIDFrom = transaction.AccountIDFrom
		output.AccountIDTo = transaction.AccountIDTo
		output.Amount = transaction.Amount
		output.CreatedAt = transaction.CreatedAt
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.TransactionFound.SetPayload(output)
	u.EventDispatcher.Dispatch(u.TransactionFound)

	return output, nil
}

func (u *FindTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := u.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.TransactionGateway)
}
