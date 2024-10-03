package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
)

type BalancesDB struct {
	DB *sql.DB
}

func NewBalancesDB(db *sql.DB) *BalancesDB {
	return &BalancesDB{
		DB: db,
	}
}

func (t *BalancesDB) UpdateBalances(account_id_from string, balanceChange float64) error {
	var currentBalance float64
	err := t.DB.QueryRow("SELECT balance FROM accounts WHERE id = ?", account_id_from).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to fetch current balance for account %s: %v", account_id_from, err)
	}

	stmt, err := t.DB.Prepare("UPDATE accounts SET balance = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(balanceChange, account_id_from)
	if err != nil {
		return fmt.Errorf("failed to update balance for account %s: %v", account_id_from, err)
	}

	fmt.Printf("Updated balance for account %s: old balance = %.2f, change = %.2f, new balance = %.2f\n",
		account_id_from, currentBalance, balanceChange, balanceChange)

	return nil
}

func (t *BalancesDB) FindBalances(accountID string) (*entity.Account, error) {
	account := &entity.Account{}

	stmt, err := t.DB.Prepare("SELECT id, client_id, balance, created_at FROM accounts WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(accountID)
	var createdAtRaw []byte
	err = row.Scan(
		&account.ID,
		&account.ClientID,
		&account.Balance,
		&createdAtRaw,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no account found with account_id: %s", accountID)
		}
		return nil, fmt.Errorf("error during scan: %v", err)
	}

	createdAtStr := string(createdAtRaw)
	account.CreatedAt, err = time.Parse("2006-01-02", createdAtStr) // assuming the format is DATE (YYYY-MM-DD)
	if err != nil {
		fmt.Println("Error parsing created_at:", err)
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}

	fmt.Println("Account balance found:", account)
	return account, nil
}
