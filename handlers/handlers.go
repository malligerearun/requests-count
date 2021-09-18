package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

var rc  = NewRequestCount()

func API() http.Handler {
	appMux := httptreemux.NewContextMux()
	appMux.Handle(http.MethodGet, "/requests-count", rc.requestsCount)
	return appMux
}