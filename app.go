package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// TranslationItem .
type TranslationItem struct {
	ID          string `db:"id"`
	Translation string `db:"translation"`
	Source      string `db:"source"`
	Count       int    `db:"count"`
	Lang        string `db:"lang"`
}

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
	translationItem := TranslationItem{}

	err = db.Get(&translationItem, "select * from vocabulary where source = $1", args)

	if err != nil {
		log.Fatalln(err)
		translation := run(args)
		db.MustExec("insert into vocabulary (source, translation, lang) values ($1, $2, $3)", args, translation, lang)
	} else {
		println(translationItem.Translation)
	}
}
