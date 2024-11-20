package main

import (
	"fmt"
	"github.com/jimbot9k/opiso/internal/error"
	"github.com/jimbot9k/opiso/internal/reverse"
	"github.com/jimbot9k/opiso/internal/cors"
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

	clientHost, clientHostFound := os.LookupEnv("FRONTEND_HOST")
	if !clientHostFound {
		clientHost = fmt.Sprintf("http://127.0.0.1:%s", port)
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /reverse", reverse.ReverseHandler)
	router.HandleFunc("/", error.NotFoundHandler)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      cors.CorsMiddleware(router, clientHost),
	}

	log.Printf("CORS Allowed for %s", clientHost)
	log.Printf("Opiso Listening on http://127.0.0.1:%s", port)
	log.Fatal(s.ListenAndServe())
}
