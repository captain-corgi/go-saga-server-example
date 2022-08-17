package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/captain-corgi/go-saga-server-example/cmd/playground/handler"
)

var (
	// store simulate database
	store = make([]string, 0)
	// storeCom simulate compensation database
	storeCom = make([]string, 0)
)

// App constants
const (
	host = ":8080"
)

// API constants
const (
	health       = "/health"
	compensation = "/compensation"
	add          = "/add"
	addTxn       = "/addTxn"
	remove       = "/remove"
	removeTxn    = "/removeTxn"
	check        = "/check"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc(health, handler.Health())

	// Routes consist of a path and a handler function.
	AddRouter(r)

	rCom := r.PathPrefix(compensation).Subrouter()
	AddCompRouter(rCom)

	// Bind to a port and pass our router in
	fmt.Printf("Http Server hanlding %s\n", host)
	log.Fatal(http.ListenAndServe(host, r))
}

func AddRouter(r *mux.Router) {
	r.HandleFunc(add, handler.Add(&store))
	r.HandleFunc(addTxn, handler.AddTxn())
	r.HandleFunc(remove, handler.Remove(&store, &storeCom))
	r.HandleFunc(removeTxn, handler.RemoveTxn())
	r.HandleFunc(check, handler.Check(&store))
}

// AddCompRouter register all compensation APIs
func AddCompRouter(r *mux.Router) {
	r.HandleFunc(add, handler.AddCom(&store))
	r.HandleFunc(remove, handler.RemoveCom(&store, &storeCom))
}
