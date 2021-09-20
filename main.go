package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/malligerearun/requests-count/handlers"
)

func main() {
	requestsCount := make(chan map[int64]int32)

	apiServer := http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: handlers.API(requestsCount),
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main: API listening on %s", apiServer.Addr)
		serverErrors <- apiServer.ListenAndServe()
	}()

	var count map[int64]int32
	go func() {
		for rc := range requestsCount {
			count = rc
		}
		close(requestsCount)
	}()

	select {
	case err := <-serverErrors:
		log.Printf("main:: server error %v", err)

	case sig := <-shutdown:
		log.Printf("main: %v: Start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := apiServer.Shutdown(ctx); err != nil {
			apiServer.Close()
			log.Printf("main::could not stop server gracefully %v", err)
		}

		f, err := os.Create("requestCountFile")
		if err != nil {
			log.Printf("main::could not create requestCountFile %v", err)
		}

		encoder := gob.NewEncoder(f)
		err = encoder.Encode(count)
		if err != nil {
			log.Printf("main::could not encode data to requestCountFile %v", err)
		}
		log.Printf("main: %v: Completed shutdown", sig)
	}
}
