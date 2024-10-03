package updatebalances

import (
	"context"
	"log"
	"time"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/gateway"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/uow"
)

type KafkaMessage struct {
	Name    string                `json:"Name"`
	Payload UpdateBalanceInputDTO `json:"Payload"`
}

type UpdateBalanceInputDTO struct {
	AccountIDFrom      string  `json:"account_id_from"`
	AccountIDTo        string  `json:"account_id_to"`
	BalanceAccountFrom float64 `json:"balance_account_id_from"`
	BalanceAccountTo   float64 `json:"balance_account_id_to"`
}
type UpdateBalanceOutputDTO struct {
	ID            string    `json:"id"`
	AccountIDFrom string    `json:"account_id_from"`
	AccountIDTo   string    `json:"account_id_to"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type UpdateBalanceUseCase struct {
	Uow             uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	BalanceUpdated  events.EventInterface
}

func NewUpdateBalanceUseCase(uow uow.UowInterface, eventDispatcher events.EventDispatcherInterface, balanceUpdated events.EventInterface) *UpdateBalanceUseCase {
	return &UpdateBalanceUseCase{
		Uow:             uow,
		EventDispatcher: eventDispatcher,
		BalanceUpdated:  balanceUpdated,
	}
}

func (u *UpdateBalanceUseCase) Execute(ctx context.Context, input UpdateBalanceInputDTO) error {
	output := &UpdateBalanceInputDTO{}
	err := u.Uow.Do(ctx, func(_ *uow.Uow) error {
		transactionRepository := u.getTransactionRepository(ctx)

		err := transactionRepository.UpdateBalances(input.AccountIDFrom, input.BalanceAccountFrom)
		if err != nil {
			log.Fatal("Error until register new transaction", err)
		}

		err = transactionRepository.UpdateBalances(input.AccountIDTo, input.BalanceAccountTo)
		if err != nil {
			log.Fatal("Error until register new transaction", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	u.BalanceUpdated.SetPayload(output)
	u.EventDispatcher.Dispatch(u.BalanceUpdated)

	return nil
}

func (u *UpdateBalanceUseCase) getTransactionRepository(ctx context.Context) gateway.BalancesGateway {
	repo, err := u.Uow.GetRepository(ctx, "BalancesDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.BalancesGateway)
}
