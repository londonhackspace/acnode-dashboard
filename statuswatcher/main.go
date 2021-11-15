package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/londonhackspace/acnode-dashboard/apiclient"
	"github.com/londonhackspace/acnode-dashboard/statuswatcher/slackapi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	consoleLogger := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(consoleLogger)
	ctx, cancel := context.WithCancel(context.Background())

	// Load the config
	cfg := GetConfig()

	var slackClient slackapi.SlackAPI = nil
	var ircCatClient *IRCCatClient = nil

	if len(cfg.SlackToken) > 0 {
		slackClient = slackapi.CreateSlackAPI(cfg.SlackToken)
	}

	if len(cfg.IRCCat) > 0 {
		ircCatClient = MakeIRCCatClient(cfg.IRCCat, cfg.IRCChannel)
	}

	// set up a signal handler to cancel the context
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Info().Str("Signal", sig.String()).Msg("Received signal")
		cancel()
	}()

	// Set up the watcher itself
	apiClient := apiclient.MakeAPIClient(cfg.ApiKey, cfg.ACNodeDashAPI)
	watcher := CreateStatusWatcher(apiClient, slackClient, cfg.SlackChannel, ircCatClient)
	go watcher.Run(ctx)

	// Set up a webserver for stats and stuff
	rtr := mux.NewRouter()
	rtr.Handle("/metrics", promhttp.Handler())

	ws := http.Server{
		Addr: cfg.HTTPListen,
		Handler: rtr,
	}

	// run the webserver
	go func() {
		log.Info().Str("ListenAddress", ws.Addr).Msg("Starting webserver")
		err := ws.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Err(err).Msg("Error in HTTP server")
		}
	}()

	// wait for the context to be cancelled
	<- ctx.Done()

	// shutdown the server
	log.Info().Msg("Shutting down server")
	timeoutctx,_ := context.WithTimeout(context.Background(), time.Second*30)
	ws.Shutdown(timeoutctx)

	log.Info().Msg("Exiting")
}
