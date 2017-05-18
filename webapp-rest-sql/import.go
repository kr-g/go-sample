package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	r := csv.NewReader(csvInput)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading all lines: %v", err)
	}

	// remove the db if existing
	os.Remove("./sample.db")

	// create sample db
	db, err := sql.Open("sqlite3", "./sample.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table tlds (id char[36] not null primary key, name text, tld text, count int default 0 );
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	for i, line := range lines {
		if i == 0 {
			// dont process header line
			continue
		}

		// create a uniq id
		uid := uuid.New().String()

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into tlds(id, name, tld) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(uid, line[0], line[1])
		if err != nil {
			log.Fatal(err)
		}

		tx.Commit()

		fmt.Println("inserted", uid, line[0], line[1])
	}

}

var csvInput = strings.NewReader(`country,tld
Switzerland,ch
Germany,de
Austria,at
France,fr
Italy,it
Spain,es
Portugal,pt
United Kingdom,uk
Ireland,ie
`)
