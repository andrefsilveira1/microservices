package createtransaction

import (
	"context"

	"github.com/andrefsilveira1/microservices/internal/entity"
	"github.com/andrefsilveira1/microservices/internal/gateway"
	"github.com/andrefsilveira1/microservices/pkg/events"
	"github.com/andrefsilveira1/microservices/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID            string
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionUseCase struct {
	Uow                uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

func NewCreateTransactionUseCase(Uow uow.UowInterface, eventDispatcher events.EventDispatcherInterface, transactionCreated events.EventInterface, balanceUpdated events.EventInterface) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow:                Uow,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated:     balanceUpdated,
	}
}

func (u *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	balanceUpdatedOutput := &BalanceUpdatedOutputDTO{}
	err := u.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := u.getAccountRepository(ctx)
		transactionRepository := u.getTransactionRepository(ctx)

		accountFrom, err := accountRepository.Find(input.AccountIDFrom)
		if err != nil {
			return err
		}

		accountTo, err := accountRepository.Find(input.AccountIDTo)
		if err != nil {
			return err
		}

		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}

		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}

		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount
		balanceUpdatedOutput.AccountIDFrom = input.AccountIDFrom
		balanceUpdatedOutput.AccountIDTo = input.AccountIDTo
		balanceUpdatedOutput.BalanceAccountIDFrom = accountFrom.Balance
		balanceUpdatedOutput.BalanceAccountIDTo = accountTo.Balance
		return nil

	})

	if err != nil {
		return nil, err
	}
	u.TransactionCreated.SetPayload(output)
	u.EventDispatcher.Dispatch(u.TransactionCreated)

	u.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	u.EventDispatcher.Dispatch(u.BalanceUpdated)

	return output, nil
}

func (u *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := u.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.AccountGateway)
}

func (u *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := u.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.TransactionGateway)
}
