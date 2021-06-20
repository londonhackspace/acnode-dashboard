package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/acserver_api"
	"github.com/londonhackspace/acnode-dashboard/acserverwatcher"
	"github.com/londonhackspace/acnode-dashboard/api"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/londonhackspace/acnode-dashboard/usagelogs"
	"log/syslog"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

var error404Template *template.Template = nil
func handle404(w http.ResponseWriter, r *http.Request) {
	if error404Template == nil {
		error404Template = getTemplate("404.gohtml")
	}
	error404Template.ExecuteTemplate(w, "404.gohtml", nil)
}

var swaggerTemplate *template.Template = nil
func handleSwagger(w http.ResponseWriter, r *http.Request) {

	if swaggerTemplate == nil {
		swaggerTemplate = getTemplate("swagger.gohtml")
	}
	swaggerTemplate.ExecuteTemplate(w, "swagger.gohtml", GetBaseTemplateArgs())
}

func main() {
	// setup logging
	consoleLogger := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(consoleLogger)

	startupWg := sync.WaitGroup{}

	conf := config.GetCombinedConfig("acserverdash.json")

	if !conf.Validate() {
		log.Fatal().Msg("Invalid configuration")
		return
	}

	if conf.SyslogServer != "" {
		sls,err := syslog.Dial("tcp", conf.SyslogServer, syslog.LOG_WARNING | syslog.LOG_DAEMON, "")
		if err == nil {
			logCombo := zerolog.MultiLevelWriter(consoleLogger, zerolog.SyslogLevelWriter(sls))
			log.Logger = log.Output(logCombo)
		}
	}

	acnodehandler := acnode.CreateACNodeHandler()

	websockerServer := api.CreateWebsockerHandler(&acnodehandler)

	if conf.LdapEnable {
		ldapauth := auth.GetLDAPAuthenticator(&conf)
		auth.AddProvider(&ldapauth)
	}

	var usageLogger usagelogs.UsageLogger = nil


	if conf.RedisEnable {
		redisConn := redis.NewClient(&redis.Options{
			Addr: conf.RedisServer,
			Password: "",
			DB: 0,
		})

		sessStore := auth.CreateRedisSessionStore(redisConn)
		auth.SetSessionStore(sessStore)
		userStore := auth.CreateRedisProvider(redisConn)
		auth.AddProvider(userStore)

		acnodehandler.SetRedis(redisConn, &startupWg)
		usageLogger = usagelogs.CreateRedisUsageLogger(redisConn)
	}

	apihandler := api.CreateApi(&conf, &acnodehandler, usageLogger)

	acserverapi := acserver_api.CreateACServer(&conf)
	acsw := acserverwatcher.Watcher{ acserverapi, &acnodehandler }

	mqttHandler := CreateMQTTHandler(&conf, &acnodehandler, usageLogger)

	// Don't initialise MQTT or ACServer watcher until everything else is running
	go func() {
		startupWg.Wait()
		go acsw.Run()
		mqttHandler.Init()
	}()


	// create a URL router
	rtr := mux.NewRouter()
	rtr.NotFoundHandler = http.HandlerFunc(handle404)

	rtr.PathPrefix("/api/").Handler(http.StripPrefix("/api", apihandler.GetRouter()))
	rtr.PathPrefix("/ws").Handler(websockerServer)
	
	// Cache the assets, unless there's no version, in which case
	// it's most likely a development version
	staticCachePolicy := CachePolicyAlways
	if getVersion() == "Unknown" {
		staticCachePolicy = CachePolicyNever
	}

	var fs http.Handler
	// if the frontend build exists, serve it from there, otherwise from /static
	if _, err := os.Stat("frontend/dist/"); !os.IsNotExist(err) {
		fs = http.FileServer(http.Dir("./frontend/dist"))
	} else {
		// serve our /static
		fs = http.FileServer(http.Dir("./static/"))
	}
	// always serve swagger from /static but without cache
	swaggerfs := CreateCacheHeaderInserter(
		http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))), CachePolicyNever)
	rtr.PathPrefix("/static/swagger/").Handler(swaggerfs)
	rtr.Handle("/static/api.yaml", swaggerfs)

	// serve our /static via a cache
	staticfs := CreateCacheHeaderInserter(fs, staticCachePolicy)
	rtr.PathPrefix("/static/").Handler(staticfs)

	//favicon and index don't get cache headers so we can change them
	rtr.Handle("/favicon.png", fs)
	rtr.Handle("/", fs)

	// Add Swagger for API docs
	rtr.HandleFunc("/swagger/", handleSwagger)

	// Prometheus-format metrics
	rtr.Handle("/metrics", promhttp.Handler())

	listen, ok := os.LookupEnv("LISTEN_ADDR")

	if !ok {
		listen = "localhost:8080"
	}

	log.Info().Msg("Listening on " + listen)
	handler := CreateCacheHeaderInserter(rtr, CachePolicyNever)
	http.ListenAndServe(listen, handlers.ProxyHeaders(LoggingHandler{ next: handler }))
}
