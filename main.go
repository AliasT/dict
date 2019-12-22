package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

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

// StateItem 统计条目
type StateItem struct {
	Source     string  `db:"source"`
	Percentage float32 `db:"percentage"`
	Count      int     `db:"count"`
}

// PadRight .
func PadRight(str string, max int) string {
	t := 0
	j := 0
	for _, rune := range str {
		cur := t + utf8.RuneLen(rune)
		if cur > max {
			return str[:j]
		}
		j = j + 1
		t = cur
	}

	return str
}

func output(items []StateItem) {
	// output
	fmt.Printf("%-20s%-10s%s\n\n", "Word", "Count", "Percentage")

	for _, item := range items {
		fmt.Printf("%-20s%-10d%.1f%s\n", item.Source, item.Count, item.Percentage, "%")
	}
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

		err = db.Get(
			&translationItem,
			"select * from vocabulary where source = $1 and lang = $2",
			args,
			lang,
		)

		if err != nil {
			// panic(err.Error())

			translation := run(":"+lang, args)
			println(translation)
			db.MustExec(
				"insert into vocabulary (source, translation, lang) values ($1, $2, $3)",
				args, translation, lang,
			)
		} else {
			// 更新查询次数
			db.MustExec(
				"update vocabulary set count = count + 1 where id = $1",
				translationItem.ID,
			)
			println(translationItem.Translation)
		}

	case "list":
		items := make([]StateItem, 0)
		err = db.Select(&items, `
			select
				source,
				count,
				count * 100.0 / (select sum(count) from vocabulary) as percentage
			from vocabulary group by source;
		`)
		if err != nil {
			panic(err.Error())
		}

		output(items)

	case "most":
		items := make([]StateItem, 0)
		err = db.Select(&items, `
			select
			source, count, count * 100.0 / (select sum(count) from vocabulary) as percentage
			from vocabulary group by source having  count >= 2;
		`)
		if err != nil {
			panic(err.Error())
		}
		output(items)
	default:
		log.Fatalln(fmt.Sprintf("Command %s not supported ", os.Args[1]))
	}
}
