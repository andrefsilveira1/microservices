package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andrefsilveira1/microservices/internal/database"
	"github.com/andrefsilveira1/microservices/internal/database/seed"
	"github.com/andrefsilveira1/microservices/internal/event"
	"github.com/andrefsilveira1/microservices/internal/event/handler"
	createaccount "github.com/andrefsilveira1/microservices/internal/usecase/create_account"
	createclient "github.com/andrefsilveira1/microservices/internal/usecase/create_client"
	createtransaction "github.com/andrefsilveira1/microservices/internal/usecase/create_transaction"
	"github.com/andrefsilveira1/microservices/internal/web"
	"github.com/andrefsilveira1/microservices/internal/web/server"
	"github.com/andrefsilveira1/microservices/pkg/events"
	"github.com/andrefsilveira1/microservices/pkg/kafka"
	"github.com/andrefsilveira1/microservices/pkg/uow"
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
	time.Sleep(5 * time.Second)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()
	seed.DropTables(db)
	seed.CreateTables(db)

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewBalanceUpdatedKafkaHandler(kafkProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDb(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransacstionUseCase := createtransaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	server := server.NewServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransacstionUseCase)

	server.AddHandler("/clients", clientHandler.CreateClient)
	server.AddHandler("/accounts", accountHandler.CreateAccount)
	server.AddHandler("/transactions", transactionHandler.CreateTransaction)

	server.Start()
}
