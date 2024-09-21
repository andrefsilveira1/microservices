package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	client1, _ := NewClient("andre", "andre.com")
	client2, _ := NewClient("freitas", "freitas.com")

	account1 := NewAccount(client1)
	account2 := NewAccount(client2)

	account1.Credit(200)
	account2.Credit(200)

	transaction, err := NewTransaction(account1, account2, 100)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 300.0, account2.Balance)
	assert.Equal(t, 100.0, account1.Balance)
}

func TestCreateTransactionWithInsufficientMoney(t *testing.T) {
	client1, _ := NewClient("andre", "andre.com")
	client2, _ := NewClient("freitas", "freitas.com")

	account1 := NewAccount(client1)
	account2 := NewAccount(client2)

	account1.Credit(200)
	account2.Credit(200)

	transaction, err := NewTransaction(account1, account2, 2000)
	assert.NotNil(t, err)
	assert.Error(t, err, "insufficient money")
	assert.Nil(t, transaction)
}
