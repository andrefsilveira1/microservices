package event

import "time"

type TransactionCreated struct {
	Name    string
	Payload interface{}
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
	}
}

func (c *TransactionCreated) GetName() string {
	return c.Name
}

func (c *TransactionCreated) GetPayload() interface{} {
	return c.Payload
}

func (e *TransactionCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (t *TransactionCreated) GetDate() time.Time {
	return time.Now()
}
