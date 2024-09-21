package handler

import (
	"fmt"
	"sync"

	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/kafka"
)

type TransactionFoundKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionFoundKafkaHandler(kafka *kafka.Producer) *TransactionFoundKafkaHandler {
	return &TransactionFoundKafkaHandler{
		Kafka: kafka,
	}
}

func (h *TransactionFoundKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message, nil, "transactions")
	fmt.Println("Transactions found with kafka")
}
