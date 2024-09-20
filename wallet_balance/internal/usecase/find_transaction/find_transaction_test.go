package gettransaction

import (
	"context"
	"testing"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/mocks"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Find(id string) (*entity.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func TestFindTransactionUseCase(t *testing.T) {
	inputDto := FindTransactionInputDTO{
		ID: "123456",
	}

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionFound()
	ctx := context.Background()
	usecase := NewFindTransactionUseCase(mockUow, dispatcher, eventTransaction)

	output, err := usecase.Execute(ctx, inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output)
}
