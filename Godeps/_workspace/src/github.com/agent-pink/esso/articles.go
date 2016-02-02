package esso

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Error struct {
	File   string
	Reason error
}

func (e Error) Error() string {
	return fmt.Sprintf("While processing %s: %s", e.File, e.Reason.Error())
}

type Article struct {
	Meta
	Contents string
}

func (a *Article) HtmlContents() template.HTML {
	return template.HTML(a.Contents)
}

type Meta struct {
	Title, Author, Slug string
	Posted              time.Time
}

var chicago *time.Location

func init() {
	var err error
	chicago, err = time.LoadLocation("America/Chicago")
	if err != nil {
		panic(err)
	}
}

func (m *Meta) Time() time.Time {
	return m.Posted.In(chicago)
}

type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}
func (a Articles) Less(i, j int) bool {
	return a[i].Posted.After(a[j].Posted)
}
func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func LoadArticle(name string) (*Article, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, Error{name, err}
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	hdr := make([]byte, 0, 255)
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		hdr = append(hdr, s.Bytes()...)
	}
	if s.Err() != nil {
		return nil, Error{name, s.Err()}
	}
	var meta Meta
	err = json.Unmarshal(hdr, &meta)
	if err != nil {
		return nil, Error{name, err}
	}
	contents := []byte{}
	for s.Scan() {
		contents = append(contents, s.Bytes()...)
	}
	if s.Err() != nil {
		return nil, Error{name, s.Err()}
	}
	return &Article{Meta: meta, Contents: string(contents)}, nil
}
func LoadArticles(pat string) (Articles, error) {
	articles := Articles{}
	names, err := filepath.Glob(pat)
	if err != nil {
		return nil, Error{pat, err}
	}
	for _, name := range names {
		article, err := LoadArticle(name)
		if err != nil {
			return nil, Error{name, err}
		}
		articles = append(articles, article)
	}
	sort.Sort(articles)
	return articles, nil
}

type ArticleMap map[string]*Article

func (a Articles) ArticleMap() ArticleMap {
	articleMap := ArticleMap{}
	for _, article := range a {
		articleMap[article.Slug] = article
	}
	return articleMap
}

type Page struct {
	Articles
	Title string
}
