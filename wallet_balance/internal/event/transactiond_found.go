package event

import "time"

type TransactionFound struct {
	Name    string
	Payload interface{}
}

func NewTransactionFound() *TransactionFound {
	return &TransactionFound{
		Name: "TransactionFound",
	}
}

func (c *TransactionFound) GetName() string {
	return c.Name
}

func (c *TransactionFound) GetPayload() interface{} {
	return c.Payload
}

func (e *TransactionFound) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (t *TransactionFound) GetDate() time.Time {
	return time.Now()
}
