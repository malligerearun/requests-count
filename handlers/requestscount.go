package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


func requestsCount(w http.ResponseWriter, r *http.Request)  {
	data := struct {		
		Count int `json:"Number of requests in the last 60 seconds"`
	} {
		Count: 10,
	}
	
	respond(data, w)
}

func respond(data interface{}, w http.ResponseWriter)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(data)
	if err!=nil {
		log.Println(err)
		return
	}
	if _, err := w.Write(jsonData); err!=nil {
		fmt.Println(err)
		return
	}
}