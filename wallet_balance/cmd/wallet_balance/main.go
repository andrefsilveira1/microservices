package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/database"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event/handler"
	gettransaction "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/find_transaction"
	registertransaction "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/register_transaction"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/web"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/web/server"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/kafka"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	consumer, err := ckafka.NewConsumer(&configMap)
	if err != nil {
		log.Fatalf("Error until kafka consumer creating %s", err)
	}

	defer consumer.Close()

	err = consumer.Subscribe("transactions", nil)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic: %s", err)
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionFound", handler.NewTransactionFoundKafkaHandler(kafkaProducer))
	eventFoundEvent := event.NewTransactionFound()
	eventRegistered := event.NewTransactionRegistered()

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	findTransactionUseCase := gettransaction.NewFindTransactionUseCase(uow, eventDispatcher, eventFoundEvent)
	registerTransactionUseCase := registertransaction.NewRegisterTransactionUseCase(uow, eventDispatcher, eventRegistered)
	server := server.NewServer(":8000")

	transactionHandler := web.NewWebTransactionHandler(*findTransactionUseCase)
	server.AddHandler("/transactions", transactionHandler.FindTransaction)

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Transaction received: %s \n", string(msg.Value))

				var kafkaMsg registertransaction.KafkaMessage
				if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
					log.Printf("Error unmarshalling Kafka message: %v", err)
					continue
				}
				payload := kafkaMsg.Payload
				log.Printf("Processing payload: %+v", payload)
				output, err := registerTransactionUseCase.Execute(ctx, payload)
				if err != nil {
					log.Printf("Error executing register transaction use case: %v", err)
				} else {
					log.Printf("Transaction successfully processed, Output ID: %s", output.ID)
				}
			} else {
				log.Printf("Error reading message from Kafka: %v", err)
			}
		}
	}()

	server.Start()
}
