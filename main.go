package main

import (
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/acserver_api"
	"github.com/londonhackspace/acnode-dashboard/acserverwatcher"
	"github.com/londonhackspace/acnode-dashboard/api"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/londonhackspace/acnode-dashboard/config"

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
	ok, _ := auth.CheckAuthUser(w, r)

	if !ok {
		http.Redirect(w, r, "/login?next=" + r.URL.Path, 302)
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
	err := indexTemplate.ExecuteTemplate(w, "index.gohtml", GetBaseTemplateArgs())
	if err != nil {
		log.Err(err).
			Str("Template", "index.gohtml").
			Msg("Error rendering template")
	}
}

var error404Template *template.Template = nil
func handle404(w http.ResponseWriter, r *http.Request) {
	if error404Template == nil {
		error404Template = getTemplate("404.gohtml")
	}
	error404Template.ExecuteTemplate(w, "404.gohtml", nil)
}

var swaggerTemplate *template.Template = nil
func handleSwagger(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

	if swaggerTemplate == nil {
		swaggerTemplate = getTemplate("swagger.gohtml")
	}
	swaggerTemplate.ExecuteTemplate(w, "swagger.gohtml", GetBaseTemplateArgs())
}

func main() {
	// setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf := config.GetCombinedConfig("acserverdash.json")

	if !conf.Validate() {
		log.Fatal().Msg("Invalid configuration")
		return
	}

	if conf.LdapEnable {
		ldapauth := auth.GetLDAPAuthenticator(&conf)
		auth.AddProvider(&ldapauth)
	}

	acnodehandler := acnode.CreateACNodeHandler()

	apihandler := api.CreateApi(&conf, &acnodehandler)

	acserverapi := acserver_api.CreateACServer(&conf)
	acsw := acserverwatcher.Watcher{ acserverapi, &acnodehandler }
	go acsw.Run()

	mqttHandler := CreateMQTTHandler(&conf, &acnodehandler)
	mqttHandler.Init()

	// create a URL router
	rtr := mux.NewRouter()
	rtr.NotFoundHandler = http.HandlerFunc(handle404)

	rtr.HandleFunc("/", handleIndex)
	rtr.HandleFunc("/login", handleLogin)
	rtr.HandleFunc("/logout", handleLogout)
	rtr.PathPrefix("/api/").Handler(http.StripPrefix("/api", apihandler.GetRouter()))
	
	// Cache the assets, unless there's no version, in which case
	// it's most likely a development version
	staticCachePolicy := CachePolicyAlways
	if getVersion() == "Unknown" {
		staticCachePolicy = CachePolicyNever
	}
	fs := CreateCacheHeaderInserter(http.FileServer(http.Dir("./static/")), staticCachePolicy)

	// Version the static directory so it can be cached
	rtr.PathPrefix(GetStaticPath()+"/").Handler(http.StripPrefix(GetStaticPath(), fs))

	// Add Swagger for API docs
	rtr.HandleFunc("/swagger/", handleSwagger)

	listen, ok := os.LookupEnv("LISTEN_ADDR")

	if !ok {
		listen = "localhost:8080"
	}

	log.Info().Msg("Listening on " + listen)
	handler := CreateCacheHeaderInserter(rtr, CachePolicyNever)
	http.ListenAndServe(listen, LoggingHandler{ next: handler })
}
