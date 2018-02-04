package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

const (
	templateFolder = "templates"
)

// Represents single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// Handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		fullpath := filepath.Join(templateFolder, t.filename)
		t.templ = template.Must(template.ParseFiles(fullpath))
	})

	t.templ.Execute(w, nil)
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
