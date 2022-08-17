package handler

import (
	"fmt"
	"log"
	"net/http"
)

// Check function return current store info
func Check(store *[]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("List: %s", *store)
		w.Write([]byte(fmt.Sprintf("List: %s", *store)))
	}
}
