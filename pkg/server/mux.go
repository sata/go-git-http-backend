package server

import "net/http"

type Router interface {
	http.Handler

	// HandleFunc registers given handler for a given pattern
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}
