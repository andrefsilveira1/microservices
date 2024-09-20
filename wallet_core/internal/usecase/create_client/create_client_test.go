package createclient

import (
	"testing"

	"github.com/andrefsilveira1/microservices/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Add(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func TestCreateClientUseCase(t *testing.T) {
	m := &ClientGatewayMock{}
	m.On("Add", mock.Anything).Return(nil)

	usecase := NewCreateClientUseCase(m)
	output, err := usecase.Execute(CreateClientInputDTO{
		Name:  "Andre",
		Email: "andre.com",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "Andre", output.Name)
	assert.Equal(t, "andre.com", output.Email)
	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Add", 1)

}
