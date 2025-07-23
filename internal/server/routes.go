package server

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", IndexHandler)
	// post request handler for GitHub webhook
	mux.HandleFunc("/github-webhook", SendNotification)

	return mux
}