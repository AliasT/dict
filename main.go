package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
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

// output
func output(items []StateItem) {
	table := tablewriter.NewWriter(os.Stdout)
	for _, item := range items {
		table.Append([]string{item.Source, fmt.Sprintf("%d", item.Count), fmt.Sprintf("%.1f%%", item.Percentage)})
	}

	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Word", "Count", "Percentage"})
	table.Render()
}

func run(args ...string) string {
	// 使用前请安装google翻译控制台工具
	// https://github.com/soimort/translate-shell
	cmd := exec.Command("trans", args...)
	cmd.Stderr = os.Stderr
	output, _ := cmd.Output()
	return string(output)
}

// Dict .
type Dict struct {
	db *sqlx.DB
}

// List 列出所有记录
func (dict *Dict) List() {
	items := make([]StateItem, 0)
	err := dict.db.Select(&items, `
		select
		source,
		count,
		count * 100.0 / (select sum(count) from vocabulary) as percentage
		from vocabulary group by source limit 10;
	`)
	if err != nil {
		panic(err.Error())
	}
	output(items)
}

// Most 返回查询次数较多的单词
func (dict *Dict) Most() {
	items := make([]StateItem, 0)
	err := dict.db.Select(&items, `
		select
		source,
		count,
		count * 100.0 / (select sum(count) from vocabulary) as percentage
		from vocabulary group by source having  count >= 2 order by count desc limit 10;
	`)
	if err != nil {
		panic(err.Error())
	}
	output(items)
}

// Query .
func (dict *Dict) Query(lang, args string) {
	translationItem := TranslationItem{}

	err := dict.db.Get(
		&translationItem,
		"select * from vocabulary where source = $1 and lang = $2",
		args,
		lang,
	)

	if err != nil {
		translation := run(":"+lang, args)
		println(translation)
		dict.db.MustExec(
			"insert into vocabulary (source, translation, lang) values ($1, $2, $3)",
			args, translation, lang,
		)
	} else {
		// 更新查询次数
		dict.db.MustExec(
			"update vocabulary set count = count + 1 where id = $1",
			translationItem.ID,
		)
		println(translationItem.Translation)
	}
}

func initDict() Dict {
	// db
	gopath := os.Getenv("GOPATH")
	// https: //stackoverflow.com/questions/32649770/how-to-get-current-gopath-from-code
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	db, err := sqlx.Connect("sqlite3", filepath.Join(gopath, "./src/github.com/Aliast/dict/dict.db"))

	if err != nil {
		log.Fatalln(err)
	}

	dict := Dict{db}
	return dict
}

func main() {
	dict := initDict()
	switch os.Args[1] {
	case "zh", "en":
		args := strings.Join(os.Args[2:], " ")
		dict.Query(os.Args[1], args)
	case "list":
		dict.List()
	case "most":
		dict.Most()

	default:
		log.Fatalln(fmt.Sprintf("Command %s not supported ", os.Args[1]))
	}
}
