package handlers

import (
	"net/http"

	"github.com/malligerearun/requests-count/web"
)

var rc = NewRequestCount()

func API(requestsCount chan map[int64]int32) http.Handler {
	app := web.NewApp(requestsCount)
	app.Handle(http.MethodGet, "/requests-count", rc.requestsCount)
	return app.Mux
}
