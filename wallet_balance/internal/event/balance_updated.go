package event

import "time"

type BalanceUpdated struct {
	Name    string
	Payload interface{}
}

func NewBalanceUpdated() *BalanceUpdated {
	return &BalanceUpdated{
		Name: "BalanceUpdated",
	}
}

func (c *BalanceUpdated) GetName() string {
	return c.Name
}

func (c *BalanceUpdated) GetPayload() interface{} {
	return c.Payload
}

func (e *BalanceUpdated) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (t *BalanceUpdated) GetDate() time.Time {
	return time.Now()
}
