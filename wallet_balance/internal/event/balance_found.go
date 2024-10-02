package event

import "time"

type BalanceFound struct {
	Name    string
	Payload interface{}
}

func NewBalanceFound() *BalanceFound {
	return &BalanceFound{
		Name: "BalanceFound",
	}
}

func (c *BalanceFound) GetName() string {
	return c.Name
}

func (c *BalanceFound) GetPayload() interface{} {
	return c.Payload
}

func (e *BalanceFound) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (t *BalanceFound) GetDate() time.Time {
	return time.Now()
}
