package jsonb

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type Row struct {
	Created time.Time
	Deleted pq.NullTime
	Body    []byte

	// embedded error
	err error

	// database connection
	db *sql.DB
}

// unmarshall database row to struct
func (r Row) Unmarshal(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	err := json.Unmarshal(r.Body, v)
	if err != nil {
		return err
	}

	return nil
}

// return driver error
func (r Row) Err() error {
	return r.err
}
