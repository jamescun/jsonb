package jsonb

import (
	"database/sql"
	"encoding/json"
)

// the Table interface implements database data aquisition methods
type Table interface {
	// get a row matching serialised query
	QueryRow(query interface{}) *Row

	// get rows matching serialised query
	Query(query interface{}) (*Rows, error)

	// create/update row
	Save(row interface{}) error

	// delete row(s) matching serialised query
	Delete(query interface{}) error
}

// internal representation of a database table
type table struct {
	// name in database of table
	name string

	// database connections
	db *sql.DB
}

func (t table) QueryRow(query interface{}) *Row {
	row := Row{db: t.db}

	stmt, err := t.db.Prepare(`
		SELECT DISTINCT ON (body->>'id') created, body
		FROM ` + t.name + `
		WHERE
			body @> $1
			AND deleted IS NULL
		ORDER BY body->>'id', created DESC
		LIMIT 1;
	`)
	if err != nil {
		row.err = err
		return &row
	}
	defer stmt.Close()

	q, err := json.Marshal(query)
	if err != nil {
		row.err = err
		return &row
	}

	err = stmt.QueryRow(q).Scan(&row.Created, &row.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			row.err = ErrNotFound
		} else {
			row.err = err
		}

		return &row
	}

	return &row
}

func (t table) Query(query interface{}) (*Rows, error) {
	stmt, err := t.db.Prepare(`
		SELECT DISTINCT ON (body->>'id') created, body
		FROM ` + t.name + `
		WHERE
			body @> $1
			AND deleted IS NULL
		ORDER BY body->>'id', created DESC;
	`)
	if err != nil {
		return nil, err
	}

	q, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(q)
	if err != nil {
		return nil, err
	}

	return &Rows{rows: rows}, nil
}

func (t table) Save(row interface{}) error {
	stmt, err := t.db.Prepare(`
		INSERT INTO ` + t.name + `
		(created, body)
		VALUES(now(), $1);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	q, err := json.Marshal(row)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (t table) Delete(query interface{}) error {
	stmt, err := t.db.Prepare(`
		UPDATE ` + t.name + `
		SET deleted = now()
		WHERE body @> $1;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	q, err := json.Marshal(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func (t table) String() string {
	return "Table: '" + t.name + "'"
}
