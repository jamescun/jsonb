package jsonb

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ErrNotFound = errors.New("jsonb: not found")
)

type Client struct {
	// database connection
	db *sql.DB
}

// create a Client connected to a postgresql server
func New(dsn string) (*Client, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

// get instance of Table
func (c Client) Table(name string) Table {
	return &table{name: name, db: c.db}
}

// create table and index for 'id' column
func (c Client) CreateTable(name string) error {
	// create table
	_, err := c.db.Exec(`
		CREATE TABLE ` + name + ` (
			created TIMESTAMP NOT NULL,
			deleted TIMESTAMP,
			body JSONB NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// create index for id column
	_, err = c.db.Exec(`
		CREATE INDEX idx` + name + `id ON ` + name + ` USING gin ((body -> 'id'));
	`)
	if err != nil {
		return err
	}

	return nil
}
