package main

import (
	"github.com/londonhackspace/acnode-dashboard/auth"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"html/template"
	"net/http"
	"os"
)

func getTemplate(page string) *template.Template {
	return template.Must(template.ParseFiles(
		"templates/base.gohtml",
		"templates/"+page))
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	ok, _ := auth.CheckAuthUser()

	if !ok {
		http.Redirect(w, r, "/login", 302)
	}

	return ok
}

var indexTemplate *template.Template = nil
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

	if indexTemplate == nil {
		indexTemplate = getTemplate("index.gohtml")
	}
	indexTemplate.ExecuteTemplate(w, "index.gohtml", nil)
}

var loginTemplate *template.Template = nil
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if loginTemplate == nil {
		loginTemplate = getTemplate("login.gohtml")
	}

	loginTemplate.ExecuteTemplate(w, "login.gohtml", nil)
}

func main() {
	// setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// create a URL router
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", handleIndex)
	rtr.HandleFunc("/login", handleLogin)
	fs := http.FileServer(http.Dir("./static/"))
	rtr.PathPrefix("/static/").Handler(http.StripPrefix( "/static", fs))

	listen, ok := os.LookupEnv("LISTEN_ADDR")

	if !ok {
		listen = "localhost:8080"
	}

	log.Info().Msg("Listening on " + listen)
	http.ListenAndServe(listen, LoggingHandler{ next: rtr })
}
