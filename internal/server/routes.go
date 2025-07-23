package server

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", IndexHandler)
	// post request handler for GitHub webhook
	mux.HandleFunc("/github-webhook", SendNotification)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
			case "/":
				mux.ServeHTTP(w, r)
			case "/github-webhook":
				mux.ServeHTTP(w, r)
			default:
				http.NotFound(w, r)
		}
	})
}