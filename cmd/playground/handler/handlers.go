package handler

import (
	"log"
	"net/http"
)

const (
	dtmServer = "http://localhost:36789/api/dtmsvr"
	appURI    = "http://localhost:8080"
)

func Health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Healthy")
		w.Write([]byte("OK\n"))
	}
}
