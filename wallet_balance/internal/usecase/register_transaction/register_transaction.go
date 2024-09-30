package registertransaction

import (
	"context"
	"log"
	"time"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/gateway"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/uow"
)

type RegisterTransactionInputDTO struct {
	ID            string  `json:"id"`
	AccountIDFrom *string `json:"account_id_from"`
	AccountIDTo   *string `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type RegisterTransactionOutputDTO struct {
	ID            string    `json:"id"`
	AccountIDFrom *string   `json:"account_id_from"`
	AccountIDTo   *string   `json:"account_id_to"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type RegisterTransactionUseCase struct {
	Uow                   uow.UowInterface
	EventDispatcher       events.EventDispatcherInterface
	TransactionRegistered events.EventInterface
}

func NewRegisterTransactionUseCase(uow uow.UowInterface, eventDispatcher events.EventDispatcherInterface, transactiondRegistered events.EventInterface) *RegisterTransactionUseCase {
	return &RegisterTransactionUseCase{
		Uow:                   uow,
		EventDispatcher:       eventDispatcher,
		TransactionRegistered: transactiondRegistered,
	}
}

func (u *RegisterTransactionUseCase) Execute(ctx context.Context, input RegisterTransactionInputDTO) (*RegisterTransactionOutputDTO, error) {
	output := &RegisterTransactionOutputDTO{}
	err := u.Uow.Do(ctx, func(_ *uow.Uow) error {
		transactionRepository := u.getTransactionRepository(ctx)
		transaction, err := entity.NewTransaction(*input.AccountIDFrom, *input.AccountIDTo, input.Amount)
		if err != nil {
			log.Printf("Fatal error until transaction creation")
		}

		err = transactionRepository.Register(transaction)
		if err != nil {
			log.Fatal("Error until register new transaction")
		}
		output.ID = transaction.ID
		output.AccountIDFrom = &transaction.AccountIDFrom
		output.AccountIDTo = &transaction.AccountIDTo
		output.Amount = transaction.Amount
		output.CreatedAt = transaction.CreatedAt
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.TransactionRegistered.SetPayload(output)
	u.EventDispatcher.Dispatch(u.TransactionRegistered)

	return output, nil
}

func (u *RegisterTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := u.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.TransactionGateway)
}
