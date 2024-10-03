package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/database"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event/handler"
	findbalances "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/find_balances"
	updatebalances "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/update_balances"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/web"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/web/server"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/events"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/kafka"
	"github.com/andrefsilveira1/microservices/wallet_balance/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi"
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

	err = consumer.Subscribe("balances", nil)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic: %s", err)
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionFound", handler.NewTransactionFoundKafkaHandler(kafkaProducer))
	eventFoundBalance := event.NewBalanceFound()
	// eventFoundEvent := event.NewTransactionFound()
	// eventRegistered := event.NewTransactionRegistered()
	eventUpdatedBalance := event.NewBalanceUpdated()

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	// registerTransactionUseCase := registertransaction.NewRegisterTransactionUseCase(uow, eventDispatcher, eventRegistered)
	// findTransactionUseCase := gettransaction.NewFindTransactionUseCase(uow, eventDispatcher, eventFoundEvent)
	updateBalancesUseCase := updatebalances.NewUpdateBalanceUseCase(uow, eventDispatcher, eventUpdatedBalance)
	findBalancesUseCase := findbalances.NewFindBalancesUseCase(uow, eventDispatcher, eventFoundBalance)
	server := server.NewServer(":8000")

	// transactionHandler := web.NewWebTransactionHandler(*findTransactionUseCase)
	// server.AddHandler("/transactions", transactionHandler.FindTransaction)

	balancesHandler := web.NewWebBalanceHandler(*findBalancesUseCase)
	server.AddHandler(http.MethodGet, "/balances/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		input := findbalances.FindBalancesInputDTO{
			ID: id,
		}

		balance, err := balancesHandler.FindBalanceByID(input, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Respond with JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(balance)
	})

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Transaction received: %s \n", string(msg.Value))

				var kafkaMsg updatebalances.KafkaMessage
				if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
					log.Printf("Error unmarshalling Kafka message: %v", err)
					continue
				}
				payload := kafkaMsg.Payload
				log.Printf("Processing payload: %+v", payload)
				err := updateBalancesUseCase.Execute(ctx, payload)
				if err != nil {
					log.Printf("Error executing register transaction use case: %v", err)
				} else {
					log.Printf("Transaction successfully processed")
				}
			} else {
				log.Printf("Error reading message from Kafka: %v", err)
			}
		}
	}()

	server.Start()
}
