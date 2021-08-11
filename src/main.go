package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const PAGES_DIR_PATH = "../pages"
const TEMPLATES_DIR_PATH = "../templates"

type Page struct {
	Title string
	Body  []byte
}

type ErrorData struct {
	Message string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(PAGES_DIR_PATH+filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	path := filepath.Join(PAGES_DIR_PATH, filename)
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func createDir() error {
	return os.MkdirAll(PAGES_DIR_PATH, os.ModePerm)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		renderTemplate(w, "error", p)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(filepath.Join(TEMPLATES_DIR_PATH, tmpl+".html"))
	t.Execute(w, p)
}

func main() {
	err := createDir()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
