package main

import (
	"fmt"
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
	// 使用前请安装google翻译控制台工具
	// https://github.com/soimort/translate-shell
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
	case "zh", "en":
		lang = os.Args[1]
		args := strings.Join(os.Args[2:], " ")
		translationItem := TranslationItem{}

		err = db.Get(&translationItem, "select * from vocabulary where source = $1", args)

		if err != nil {
			// panic(err.Error())
			translation := run(args)
			println(translation)
			db.MustExec("insert into vocabulary (source, translation, lang) values ($1, $2, $3)", args, translation, lang)
		} else {
			// 更新查询次数
			db.MustExec("update vocabulary set count = count + 1 where id = $1", translationItem.ID)
			println(translationItem.Translation)
		}
	case "list":
		// TODO

	default:
		log.Fatalln(fmt.Sprintf("Command %s not supported ", os.Args[1]))
	}
}
