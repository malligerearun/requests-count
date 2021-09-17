package main

import (
	"fmt"
	"net/http"

	"github.com/malligerearun/requests-count/handlers"
)



func main()  {
	apiServer := http.Server {	
		Addr: "0.0.0.0:3000",
		Handler: handlers.API(),
	}
	fmt.Println("Started server ...")
	apiServer.ListenAndServe()
}