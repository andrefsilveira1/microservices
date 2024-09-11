package createtransaction

import (
	"context"
	"fmt"
	"testing"

	"github.com/andrefsilveira1/microservices/internal/entity"
	"github.com/andrefsilveira1/microservices/internal/event"
	"github.com/andrefsilveira1/microservices/internal/usecase/mocks"
	"github.com/andrefsilveira1/microservices/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transasction *entity.Transaction) error {
	args := m.Called(transasction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) Find(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	fmt.Println("")
	return nil
}

func TestCreateTransactionuseCase(t *testing.T) {
	client1, _ := entity.NewClient("andre1", "andre1.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("andre2", "andre2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        500,
	}

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	ctx := context.Background()

	usecase := NewCreateTransactionUseCase(mockUow, dispatcher, event)

	output, err := usecase.Execute(ctx, inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
