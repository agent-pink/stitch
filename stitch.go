package main

import (
	"flag"
	"github.com/agent-pink/esso"
	"html/template"
	"os"
)

func main() {
	pat := flag.String("pat", "articles/*.html", "Pattern for articles")
	title := flag.String("title", "Lorem Ipsum", "Title for page")
	flag.Parse()
	articles, err := esso.LoadArticles(*pat)
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.ParseFiles(flag.Args()...))
	tpl.Execute(os.Stdout, esso.Page{Articles: articles, Title: *title})
}
