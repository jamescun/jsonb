package jsonb

import (
	"database/sql"
)

type Rows struct {
	rows *sql.Rows
}

func (r *Rows) Unmarshal(v interface{}) error {
	var row Row
	if err := r.rows.Scan(&row.Created, &row.Body); err != nil {
		return err
	}

	return row.Unmarshal(v)
}

func (r *Rows) Close() error {
	return r.rows.Close()
}

func (r *Rows) Next() bool {
	return r.rows.Next()
}

func (r *Rows) Err() error {
	return r.rows.Err()
}
