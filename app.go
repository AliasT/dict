package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Connect("sqlite3", "dict.db")
	if err != nil {
		log.Fatalln(err)
	}

	print(db)
}
