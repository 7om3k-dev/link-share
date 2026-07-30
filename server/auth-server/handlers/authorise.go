package handlers

import (
	"crypto/rand"
	"net/http"
	"slices"
)

func Authorise(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientId := q.Get("client_id")

	client := getClientData(clientId)
	if client == nil {
		// TODO: Log error to file as someone tried to pretend to be client
		http.Error(w, "Unknown client", http.StatusBadRequest)
		return
	}

	redirectUri := q.Get("redirect_uri")
	if slices.Contains(client.redirectUris, redirectUri) == false {
		// TODO: Log error to file as someone tried to pretend to be client
		http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rand.Text()))
}

func Approve(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Endpoint not ready", http.StatusInternalServerError)
}
