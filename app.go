package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func run(args ...string) {
	add := exec.Command("git", args...)
	add.Stdout = os.Stdout
	add.Stderr = os.Stderr
	add.Run()
}

func main() {
	dictCmd := flag.NewFlagSet("mode", flag.ExitOnError)

	switch os.Args[1] {
	case "zh":
		dictCmd.Parse(os.Args[2:])
		println("zh")
	case "en":
		dictCmd.Parse(os.Args[2:])
		println("en")
	}
	// db
	db, err := sqlx.Connect("sqlite3", "dict.db")
	if err != nil {
		log.Fatalln(err)
	}

	print(db)

	if err != nil {
		log.Fatal(err)
	}
}
