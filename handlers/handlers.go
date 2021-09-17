package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)


func API() http.Handler {
	appMux := httptreemux.NewContextMux()
	appMux.Handle(http.MethodGet, "/requests-count", requestsCount)
	return appMux
}