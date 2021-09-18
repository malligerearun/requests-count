package web

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(w http.ResponseWriter, r *http.Request, requestsCount chan map[int64]int32)

type App struct {
	Mux           *httptreemux.ContextMux
	requestsCount chan map[int64]int32
}

func NewApp(requestsCount chan map[int64]int32) *App {
	mux := httptreemux.NewContextMux()
	return &App{
		Mux:           mux,
		requestsCount: requestsCount,
	}
}

func (a *App) Handle(method string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, a.requestsCount)
	}
	a.Mux.Handle(method, path, h)
}
