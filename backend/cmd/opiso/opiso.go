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

	processesAllowedRaw, processesAllowedFound := os.LookupEnv("PROCESS_COUNT")
	if !processesAllowedFound {
		processesAllowedRaw = "1000";
	}
	processesAllowed, err := strconv.Atoi(processesAllowedRaw)
    if err != nil {
		log.Fatal("Invalid Process Count Provided")
		return;
    }

	processesAllowedSemaphore := make(chan struct{}, processesAllowed)
	router := http.NewServeMux()
	router.HandleFunc("POST /reverse", reverse.ReverseHandler(processesAllowedSemaphore))
	router.HandleFunc("/", error.NotFoundHandler)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      cors.CorsMiddleware(router, corsOrigin),
	}

	log.Printf("%d Handler Processes Allowed Concurrently", processesAllowed)
	log.Printf("CORS Allowed for %s", corsOrigin)
	log.Printf("Opiso Listening on http://127.0.0.1:%s", port)
	log.Fatal(s.ListenAndServe())
}
