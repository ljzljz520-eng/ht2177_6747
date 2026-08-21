package main

import (
	"log"
	"net/http"
	"os"

	"storeledger/internal/config"
	"storeledger/internal/httpapi"
	"storeledger/internal/service"
	"storeledger/internal/store"
)

func main() {
	c := config.Load()
	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}
	st, err := store.Open(c.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()
	svc, err := service.New(st, service.DeterministicClock("1970-01-01T00:00:00Z"), c.Reviewer)
	if err != nil {
		log.Fatal(err)
	}
	server := httpapi.NewServer(svc)
	log.Fatal(http.ListenAndServe(c.HTTPAddress, server.Handler()))
}

var _ = os.Args
