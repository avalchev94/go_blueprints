package main

import (
	"flag"
	"github.com/avalchev94/go_blueprints/trace"
	"net/http"
	"os"
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

	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	r.tracer.Trace("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		r.tracer.Trace("ListenAndServe:", err)
	}
}
