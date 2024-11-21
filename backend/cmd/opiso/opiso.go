package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jimbot9k/opiso/internal/cors"
	"github.com/jimbot9k/opiso/internal/error"
	"github.com/jimbot9k/opiso/internal/reverse"
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

	corsOrigin, corsOriginFound := os.LookupEnv("CORS_ORIGIN")
	if !corsOriginFound {
		corsOrigin = fmt.Sprintf("http://127.0.0.1:%s", port)
	}

	routinesAllowedRaw, routinesAllowedFound := os.LookupEnv("ROUTINES_LIMIT")
	if !routinesAllowedFound {
		routinesAllowedRaw = "50000";
	}
	routinesAllowed, err := strconv.Atoi(routinesAllowedRaw)
    if err != nil {
		log.Fatal("Invalid Routine Count Provided")
		return;
    }

	routinesAllowedSemaphore := make(chan struct{}, routinesAllowed)
	router := http.NewServeMux()
	router.HandleFunc("POST /reverse", reverse.ReverseHandler(routinesAllowedSemaphore))
	router.HandleFunc("/", error.NotFoundHandler)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      cors.CorsMiddleware(router, corsOrigin),
	}

	log.Printf("%d Handler Routines Allowed Concurrently", routinesAllowed)
	log.Printf("CORS Allowed for %s", corsOrigin)
	log.Printf("Opiso Listening on http://127.0.0.1:%s", port)
	log.Fatal(s.ListenAndServe())
}
