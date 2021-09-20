package handlers

import (
	"net/http"

	"github.com/malligerearun/requests-count/web"
)

var rc = NewRequestCount()

func API(requestsCount chan map[int64]int32) http.Handler {
	app := web.NewApp(requestsCount)

	tr := NewTestReq(rc)

	app.Handle(http.MethodGet, "/requests-count", rc.requestsCount)
	app.Handle(http.MethodGet, "/test-request", tr.testRequest)
	return app.Mux
}
