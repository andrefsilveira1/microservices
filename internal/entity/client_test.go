package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("André Freitas", "andre.com")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "André Freitas", client.Name)
	assert.Equal(t, "andre.com", client.Email)

}

func TestCreateNewClientWhenInvalidParameters(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("Andre", "andre.com")
	err := client.Update("andre", "freitas.com")
	assert.Nil(t, err)
	assert.Equal(t, "andre", client.Name)
	assert.Equal(t, "freitas.com", client.Email)
}

func TestUpdateClientWithInvalidParameters(t *testing.T) {
	client, _ := NewClient("Andre", "andre.com")
	err := client.Update("", "freitas.com")
	assert.Error(t, err, "name is required")
	assert.Equal(t, "freitas.com", client.Email)
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("Andre", "andre.com")
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))

}
