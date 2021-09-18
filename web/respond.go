package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func Respond(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("web::respond:: error marshaling response %v", err)
		return
	}
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("web::respond:: error writing response %v", err)
		return
	}
}
