package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func run(args ...string) {
	git := exec.Command("trans", args...)
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr
	git.Run()
}

func main() {
	dictCmd := flag.NewFlagSet("mode", flag.ExitOnError)

	switch os.Args[1] {
	case "zh":
		args := os.Args[2:]
		run(strings.Join(args, " "))
		dictCmd.Parse(os.Args[2:])
	case "en":
		args := os.Args[2:]
		run(strings.Join(args, " "))
		dictCmd.Parse(os.Args[2:])
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
