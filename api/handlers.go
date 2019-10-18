package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (a *API) jsonResponse(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json;utf-8")

	b, err := json.Marshal(response{Message: msg, Status: status})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal error", "status": "500"}`))
		a.log.Errorf("Failed to marshal JSON: %v", err)
		return
	}
	w.WriteHeader(status)
	w.Write(b)
}

func (a *API) jsonResponseRaw(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json;utf-8")

	b, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "internal error", "status": "500"}`))
		a.log.Errorf("Failed to marshal JSON: %v", err)
		return
	}
	w.WriteHeader(status)
	w.Write(b)
}

func (a *API) handlerHostname(w http.ResponseWriter, r *http.Request) {
	a.jsonResponse(w, a.hostname, http.StatusOK)
	a.log.ResponseLog(r)
}

func (a *API) handlerHello(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 2 {
		a.jsonResponse(w, "Hello, "+parts[2], http.StatusOK)
		return
	}

	a.jsonResponse(w, "Hello, world", http.StatusOK)
}

func (a *API) handlerStats(w http.ResponseWriter, t *http.Request) {
	res, err := a.hits.Stats()
	if err != nil {
		a.jsonResponse(w, "error", http.StatusInternalServerError)
		a.log.Errorf("Failed to read stats: %v", err)
		return
	}
	a.jsonResponseRaw(w, res, http.StatusOK)
}
