package database

import (
	"fmt"
	"time"

	"github.com/andrefsilveira1/microservices/wallet_balance/internal/entity"
)

func (t *TransactionDB) UpdateBalances(account_id_from string, balanceChange float64) error {
	var currentBalance float64
	err := t.DB.QueryRow("SELECT balance FROM accounts WHERE id = ?", account_id_from).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to fetch current balance for account %s: %v", account_id_from, err)
	}

	newBalance := currentBalance + balanceChange

	stmt, err := t.DB.Prepare("UPDATE accounts SET balance = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newBalance, account_id_from)
	if err != nil {
		return fmt.Errorf("failed to update balance for account %s: %v", account_id_from, err)
	}

	fmt.Printf("Updated balance for account %s: old balance = %.2f, change = %.2f, new balance = %.2f\n",
		account_id_from, currentBalance, balanceChange, newBalance)

	return nil
}

func (t *TransactionDB) FindBalances(id string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	stmt, err := t.DB.Prepare("SELECT id, account_id_from, account_id_to, amount, created_at from transactions_balance WHERE id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(id)
	var createdAtRaw []byte
	err = row.Scan(
		&transaction.ID,
		&transaction.AccountIDFrom,
		&transaction.AccountIDTo,
		&transaction.Amount,
		&createdAtRaw,
	)

	if err != nil {
		fmt.Println("Error during scan:", err)
		return nil, err
	}

	createdAtStr := string(createdAtRaw)
	transaction.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		fmt.Println("Error parsing created_at:", err)
		return nil, err
	}

	fmt.Println("Transaction found:", transaction)
	return transaction, nil

}
