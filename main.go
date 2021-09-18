package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/malligerearun/requests-count/handlers"
)



func main()  {
	apiServer := http.Server {	
		Addr: "0.0.0.0:3000",
		Handler: handlers.API(),
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	serverErrors := make(chan error, 1)
	
	go func() {
		log.Printf("main: API listening on %s", apiServer.Addr)
		serverErrors <- apiServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		log.Printf("main:: server error %v", err)

	case sig := <-shutdown:
		log.Printf("main: %v: Start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		if err := apiServer.Shutdown(ctx); err != nil {
			apiServer.Close()
			log.Printf("main::could not stop server gracefully %v", err)
		}

		log.Printf("main: %v: Completed shutdown", sig)
	}
}