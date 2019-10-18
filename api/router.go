package api

import (
	"net/http"
	"regexp"
)

// bootRouter boots router
func (a *API) bootRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.router)

	return mux
}

// router for the requests
func (a *API) router(w http.ResponseWriter, r *http.Request) {
	a.log.RequestLog(r)
	path := r.URL.Path
	a.hits.Incr(path)

	rHello, err := regexp.Compile("/hello/?.*")
	if err != nil {
		a.log.Fatal("Failed to compile /hello path")
	}

	switch r.Method {
	case http.MethodGet:
		switch {
		case path == "/":
			a.handlerHostname(w, r)
		case path == "/stats":
			a.handlerStats(w, r)
		case rHello.MatchString(path):
			a.handlerHello(w, r)
		default:
			a.jsonResponse(w, "not found", http.StatusNotFound)
		}
	default:
		a.jsonResponse(w, "not found", http.StatusNotFound)
	}
}
