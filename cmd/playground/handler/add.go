package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/google/uuid"
)

func Add(store *[]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Add START\n")
		defer log.Printf("Add END\n")

		*store = append(*store, fmt.Sprintf("%d", time.Now().Minute()))

		if len(*store) > 2 {
			log.Println("Add failed: store > 2")
			w.WriteHeader(409)
			w.Write([]byte("NG\n"))
			return
		}

		log.Println("Added")
		w.Write([]byte("OK\n"))
	}
}

func AddCom(store *[]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddCom START\n")
		defer log.Printf("AddCom END\n")

		lastIndex := len(*store) - 1
		if lastIndex >= 0 {
			*store = (*store)[:lastIndex]
		}

		log.Println("Add reverted")
		w.Write([]byte("OK\n"))
	}
}

// AddTxn is a usecase using SAGA
func AddTxn() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddTxn\n")
		saga := dtmcli.NewSaga(dtmServer, uuid.New().String()).
			Add(appURI+"/add", appURI+"/compensation"+"/add", "")
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
