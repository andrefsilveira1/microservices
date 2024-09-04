package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewCleint("André Freitas", "andre.com")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "André Freitas", client.Name)
	assert.Equal(t, "andre.com", client.Email)

}

func TestCreateNewClientWhenInvalidParameters(t *testing.T) {
	client, err := NewCleint("André Freitas", "andre.com")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}
