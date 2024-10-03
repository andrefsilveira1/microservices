package findbalances

import (
	"context"
	"time"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/gateway"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/uow"
)

type FindBalancesInputDTO struct {
	ID string `json:"id"`
}

type FindBalancesOutputDTO struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"client_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type FindBalancesUseCase struct {
	Uow             uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	BalancesFound   events.EventInterface
}

func NewFindBalancesUseCase(uow uow.UowInterface, eventDispatcher events.EventDispatcherInterface, balancesFound events.EventInterface) *FindBalancesUseCase {
	return &FindBalancesUseCase{
		Uow:             uow,
		EventDispatcher: eventDispatcher,
		BalancesFound:   balancesFound,
	}
}

func (u *FindBalancesUseCase) Execute(ctx context.Context, input FindBalancesInputDTO) (*FindBalancesOutputDTO, error) {
	output := &FindBalancesOutputDTO{}
	err := u.Uow.Do(ctx, func(_ *uow.Uow) error {
		transactionRepository := u.getTransactionRepository(ctx)

		account, err := transactionRepository.FindBalances(input.ID)
		if err != nil {
			return err
		}

		output.ID = account.ID
		output.ClientID = account.ClientID
		output.Balance = account.Balance
		output.CreatedAt = account.CreatedAt
		output.UpdatedAt = account.UpdatedAt
		return nil
	})

	if err != nil {
		return nil, err
	}

	u.BalancesFound.SetPayload(output)
	u.EventDispatcher.Dispatch(u.BalancesFound)

	return output, nil
}

func (u *FindBalancesUseCase) getTransactionRepository(ctx context.Context) gateway.BalancesGateway {
	repo, err := u.Uow.GetRepository(ctx, "BalancesDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.BalancesGateway)
}
