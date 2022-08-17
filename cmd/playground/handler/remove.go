package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/google/uuid"
)

// Remove function remove an item from store, and use storeCom for store original value
func Remove(store *[]string, storeCom *[]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Remove START\n")
		defer log.Printf("Remove END\n")

		lastIndex := len(*store) - 1
		if lastIndex >= 0 {
			lastValue := (*store)[lastIndex]
			*storeCom = append(*storeCom, lastValue)
			*store = (*store)[:lastIndex]

			minVal, err := strconv.Atoi(lastValue)
			if err != nil {
				msg := fmt.Sprintf("Cannot convert %s to Int\n", lastValue)
				log.Println(msg)
				w.WriteHeader(409)
				w.Write([]byte(msg))
				return
			}
			if minVal%2 == 0 {
				msg := fmt.Sprintf("Cannot remove %d\n", minVal)
				log.Println(msg)
				w.WriteHeader(409)
				w.Write([]byte(msg))
				return
			}
		}

		log.Println("Removed")
		w.Write([]byte("OK\n"))
	}
}

// RemoveTxn function remove an item from store, and use storeCom for store original value
func RemoveTxn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddTxn\n")
		saga := dtmcli.NewSaga(dtmServer, uuid.New().String()).
			Add(appURI+"/remove", appURI+"/compensation"+"/remove", "")
		saga.WaitResult = true
		if err := saga.Submit(); err != nil {
			msg := fmt.Sprintf("SAGA transaction failed: %s", err.Error())
			log.Print(msg)
			w.Write([]byte(msg))
		} else {
			msg := "SAGA transaction success"
			log.Print(msg)
			w.Write([]byte(msg))
		}
	}
}

// RemoveCom function restore an item from storeCom in case of remove item failed
func RemoveCom(store *[]string, storeCom *[]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("RemoveCom START\n")
		defer log.Printf("RemoveCom END\n")

		lastIndex := len(*store) - 1
		if lastIndex >= 0 {
			*store = append(*store, (*storeCom)[lastIndex])
			*storeCom = (*storeCom)[:lastIndex]
		}
		log.Println("Remove reverted")
		w.Write([]byte("OK\n"))
	}
}
