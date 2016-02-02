package esso

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var articleSlice Articles
var articleHash ArticleMap

func App() *mux.Router {
	var err error
	baseTpl = template.Must(template.ParseFiles("templates/base.html"))
	articleTpl = template.Must(template.Must(baseTpl.Clone()).ParseFiles("templates/article.html"))
	articlesTpl = template.Must(template.Must(baseTpl.Clone()).ParseFiles("templates/article.html"))
	app := mux.NewRouter()
	articleSlice, err = LoadArticles("articles/*.html")
	if err != nil {
		panic(err)
	}
	articleHash = articleSlice.ArticleMap()
	app.HandleFunc("/", ArticlesHandler)
	app.HandleFunc("/articles/", ArticlesHandler)
	app.HandleFunc("/articles/{slug}", ArticleHandler)
	app.Handle("/static/{page:.*}", http.FileServer(http.Dir("public")))
	return app
}

var baseTpl, articlesTpl, articleTpl *template.Template

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{Articles: articleSlice, Title: "Essocony: All Articles"}
	articleTpl.Execute(w, data)
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	article, found := articleHash[slug]
	if !found {
		http.NotFound(w, r)
		return
	}
	data := Page{Articles: Articles{article}, Title: "Essocony: " + article.Title}
	articlesTpl.Execute(w, data)
}
