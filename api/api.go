package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/skibish/docker-demo-app/hits"
	"github.com/skibish/docker-demo-app/log"
)

// API is a structure that contains API for the bot
type API struct {
	server   *http.Server
	log      *log.Logger
	hostname string
	hits     *hits.Hits
}

// New return new API instance
func New(log *log.Logger, hostname string, hit *hits.Hits) *API {
	return &API{
		log:      log,
		hostname: hostname,
		hits:     hit,
	}
}

// Start starts the API server
func (a *API) Start(port int) error {
	p := strconv.Itoa(port)
	s := &http.Server{
		Addr:    ":" + p,
		Handler: a.bootRouter(),
	}

	a.server = s

	return s.ListenAndServe()
}

// Shutdown performs graceful API shutdown
func (a *API) Shutdown() error {
	return a.server.Shutdown(context.Background())
}
