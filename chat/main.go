package main

import (
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
	go func() {
		for {
		}
	}()
	for {

	}
}
