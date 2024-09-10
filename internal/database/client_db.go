package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/andrefsilveira1/microservices/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDb(db *sql.DB) *ClientDB {
	return &ClientDB{
		DB: db,
	}
}

func (c *ClientDB) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := c.DB.Prepare("SELECT id, name, email, DATE_FORMAT(created_at, '%Y-%m-%d') FROM clients WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var createdAtStr string

	row := stmt.QueryRow(id)
	if err := row.Scan(&client.ID, &client.Name, &client.Email, &createdAtStr); err != nil {
		return nil, err
	}

	createdAt, err := time.Parse("2006-01-02", createdAtStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil, err
	}
	client.CreatedAt = createdAt

	return client, nil
}

func (c *ClientDB) Add(client *entity.Client) error {
	stmt, err := c.DB.Prepare("INSERT INTO clients (id, name ,email, created_at) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt)
	if err != nil {
		return err
	}

	return nil

}
