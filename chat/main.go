package main

import (
	"flag"
	"github.com/avalchev94/go_blueprints/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	gomniauth.SetSecurityKey("some long key")
	gomniauth.WithProviders(
		google.New(
			"928275824275-ki45gkotd8pn262od4k9jvp235nd0et8.apps.googleusercontent.com",
			"nvZoSjcc8t9yWSdXRBmUNYaF",
			"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	r.tracer.Trace("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		r.tracer.Trace("ListenAndServe:", err)
	}
}
