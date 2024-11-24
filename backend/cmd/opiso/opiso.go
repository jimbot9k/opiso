package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/jimbot9k/opiso/internal/cors"
	"github.com/jimbot9k/opiso/internal/error"
	"github.com/jimbot9k/opiso/internal/headers"
	"github.com/jimbot9k/opiso/internal/reverse"
	"github.com/jimbot9k/opiso/internal/status"
	"github.com/jimbot9k/opiso/internal/util"
)

const DEFAULT_PORT = "8080"
const DEFAULT_ROUTINES = 5000
const DEFAULT_CACHE_SIZE = 1000
const DEFAULT_MINIMUM_WORD_SIZE_CACHE = 10

/*
A server for reversing words
*/
func main() {

	log.Printf("Opiso Starting")

	port := util.GetStringEnvironmentVariable("PORT", DEFAULT_PORT, func() bool { return true })
	corsOrigin := util.GetStringEnvironmentVariable("CORS_ORIGIN", fmt.Sprintf("http://127.0.0.1:%s", port), func() bool { return true })
	routinesAllowed := util.GetPositiveIntegerEnvironmentVariable("ROUTINE_LIMIT", DEFAULT_ROUTINES, func() bool { return true })
	cacheSize := util.GetPositiveIntegerEnvironmentVariable("CACHE_SIZE", DEFAULT_CACHE_SIZE, func() bool { return true })
	cacheWordMinimumLength := util.GetPositiveIntegerEnvironmentVariable("CACHE_MINIMUM", DEFAULT_MINIMUM_WORD_SIZE_CACHE, func() bool { return true })

	routinesAllowedSemaphore := make(chan struct{}, routinesAllowed)
	router := http.NewServeMux()
	router.HandleFunc("POST /reverse", reverse.ReverseHandlerWithCache(routinesAllowedSemaphore, cacheSize, cacheWordMinimumLength))
	router.HandleFunc("GET /health", status.HealthHandler)
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/", error.NotFoundHandler)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      headers.HeaderJsonMiddleware(cors.CorsMiddleware(router, corsOrigin)),
	}

	log.Printf("%d Handler Routines Allowed Concurrently", routinesAllowed)
	log.Printf("%d Cached Messages Allowed", cacheSize)
	log.Printf("%d Length Required for Message to Cache", cacheWordMinimumLength)
	log.Printf("CORS Allowed for %s", corsOrigin)
	log.Printf("Opiso Listening on http://127.0.0.1:%s", port)
	log.Fatal(s.ListenAndServe())
}
