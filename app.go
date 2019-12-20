package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func run(args ...string) string {
	cmd := exec.Command("trans", args...)
	cmd.Stderr = os.Stderr
	output, _ := cmd.Output()
	return string(output)
}

func main() {
	// dictCmd := flag.NewFlagSet("mode", flag.ExitOnError)

	// db
	db, err := sqlx.Connect("sqlite3", "./dict.db")

	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	var lang string
	switch os.Args[1] {
	case "zh":
		lang = "zh"
		// dictCmd.Parse(os.Args[2:])
	case "en":
		lang = "en"
		// dictCmd.Parse(os.Args[2:])
	}

	args := strings.Join(os.Args[2:], " ")
	translation := run(args)
	db.MustExec("insert into vocabulary (source, translation, lang) values ($1, $2, $3)", args, translation, lang)

}
