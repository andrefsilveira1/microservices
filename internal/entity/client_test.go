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
	client, err := NewCleint("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewCleint("Andre", "andre.com")
	err := client.Update("andre", "freitas.com")
	assert.Nil(t, err)
	assert.Equal(t, "andre", client.Name)
	assert.Equal(t, "freitas.com", client.Email)
}

func TestUpdateClientWithInvalidParameters(t *testing.T) {
	client, _ := NewCleint("Andre", "andre.com")
	err := client.Update("", "freitas.com")
	assert.Error(t, err, "name is required")
	assert.Equal(t, "freitas.com", client.Email)
}
