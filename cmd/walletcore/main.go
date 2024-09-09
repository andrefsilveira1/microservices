package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/andrefsilveira1/microservices/internal/database"
	"github.com/andrefsilveira1/microservices/internal/event"
	createaccount "github.com/andrefsilveira1/microservices/internal/usecase/create_account"
	createclient "github.com/andrefsilveira1/microservices/internal/usecase/create_client"
	createtransaction "github.com/andrefsilveira1/microservices/internal/usecase/create_transaction"
	"github.com/andrefsilveira1/microservices/internal/web"
	"github.com/andrefsilveira1/microservices/internal/web/server"
	"github.com/andrefsilveira1/microservices/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Create the connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	// eventDispatcher.Register("TransactionCreated", handler)
	transactionCreatedEvent := event.NewTransactionCreated()

	clientDb := database.NewClientDb(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransacstionUseCase := createtransaction.NewCreateTransactionUseCase(transactionDb, accountDb, eventDispatcher, transactionCreatedEvent)

	server := server.NewServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransacstionUseCase)

	server.AddHandler("/clients", clientHandler.CreateClient)
	server.AddHandler("/accounts", accountHandler.CreateAccount)
	server.AddHandler("/transactions", transactionHandler.CreateTransaction)

	server.Start()
}
