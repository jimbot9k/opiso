package main

import (
	"fmt"
	"gilmour/opiso/internal/error"
	"gilmour/opiso/internal/reverse"
	"log"
	"net/http"
	"os"
	"time"
)

const DEFAULT_PORT = "8080"

/*
A server for reversing words
*/
func main() {

	log.Printf("Opiso Starting")

	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = DEFAULT_PORT
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /reverse/{word}", reverse.ReverseHandler)
	router.HandleFunc("/", error.NotFoundHandler)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	log.Printf("Opiso Listening on http://127.0.0.1:%s", port)
	log.Fatal(s.ListenAndServe())
}
