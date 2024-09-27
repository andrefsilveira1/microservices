package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/database"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/database/seed"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event"
	"github.com/andrefsilveira1/microservices/wallet_balance/internal/event/handler"
	gettransaction "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/find_transaction"
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
	seed.CreateTables(db)

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionFound", handler.NewTransactionFoundKafkaHandler(kafkaProducer))
	eventFoundEvent := event.NewTransactionFound()

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	findTransactionUseCase := gettransaction.NewFindTransactionUseCase(uow, eventDispatcher, eventFoundEvent)

	server := server.NewServer(":8000")

	transactionHandler := web.NewWebTransactionHandler(*findTransactionUseCase)
	server.AddHandler("/transactions", transactionHandler.FindTransaction)

	server.Start()
}
