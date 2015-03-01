JSONB
=====

JSONB is a wrapper around `database/sql` to expose a simple document-store for PostgreSQL using 9.4+ JSONB types.

Features:

  - CRUD
  - Versioning
  - Similar API to `database/sql`

Example
-------

package main

	import (
		"github.com/jamescun/jsonb"
	)

	type Transformer struct {
		Id        int    `json:"id,omitempty"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	}

	func main() {
		// connect JSONB client to postgresql
		client, _ := jsonb.New("dbname=demo sslmode=disable")

		// create table for transformers
		client.CreateTable("transformers")
		transformers = client.Table("transformers")

		// create a row
		optimus := Transformer{Id: 1, FirstName: "Wheel", LastName: "Jack"}
		transformers.Save(optimus)

		transformers.Save(Transformer{Id: 2, FirstName: "Rodimus", LastName: "Prime"})

		// update row
		optimus.FirstName = "Optimus"
		optimus.LastName = "Prime"
		transformers.Save(optimus)

		// fetch a row
		var rodimus Transformer
		transformers.QueryRow(Transformer{Id: 2}).Unmarshal(&rodimus)

		// fetch all rows with LastName "Prime"
		rows, err := transformers.Query(Transformer{LastName: "Prime"})
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var transformer Transformer
			if err := rows.Unmarshal(transformer); err != nil {
				panic(err)
			}
		}

		if err := rows.Err(); err != nil {
			panic(err)
		}

		// delete all rows with FirstName "Rodimus"
		transformers.Delete(Transformer{FirstName: "Rodimus"})
	}
