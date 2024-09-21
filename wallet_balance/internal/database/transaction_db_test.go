package database

import (
	"database/sql"
	"testing"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDB TransactionDB
	transaction   *entity.Transaction
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("Create table transactions (id varchar(255),  account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date)")
	s.transactionDB = *NewTransactionDB(db)
	s.transaction, _ = entity.NewTransaction("1234", "12345", 50)

}

func (s *TransactionDBTestSuite) TearDownSwuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}
