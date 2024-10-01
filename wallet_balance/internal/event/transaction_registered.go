package event

import "time"

type TransactionRegistered struct {
	Name    string
	Payload interface{}
}

func NewTransactionRegistered() *TransactionRegistered {
	return &TransactionRegistered{
		Name: "TransactionRegistered",
	}
}

func (c *TransactionRegistered) GetName() string {
	return c.Name
}

func (c *TransactionRegistered) GetPayload() interface{} {
	return c.Payload
}

func (e *TransactionRegistered) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (t *TransactionRegistered) GetDate() time.Time {
	return time.Now()
}
