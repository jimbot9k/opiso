package reverse

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jimbot9k/opiso/internal/error"
	"github.com/prometheus/client_golang/prometheus"
)

type reverseRequest struct {
	Messages []string `json:"messages"`
}

type reverseResponse struct {
	Reversed []string `json:"reversed"`
}

var (
	semaphoreLimitUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "opiso_reverse_routine_limit_semaphore_slots_in_use",
			Help: "Number of semaphore slots currently in by reverse handler",
		},
	)
	messagesReversed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "opiso_reverse_messages_reversed",
			Help: "Number of messages reversed",
		},
	)
	totalReverseRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "opiso_reverse_reverse_requests",
			Help: "Number of requests to reverse messages",
		},
	)
	messagesCached = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "opiso_reverse_messages_cached",
			Help: "Number of messages in cache",
		},
	)
)

func init() {
	prometheus.MustRegister(semaphoreLimitUsage, messagesReversed, totalReverseRequests, messagesCached)
}

// ReverseHandlerWithCache creates an HTTP handler for reversing words with caching.
func ReverseHandlerWithCache(routinesCountSemaphore chan struct{}, cacheSize int, cacheKeyMinimumSize int) http.HandlerFunc {
	cache := NewCache(cacheSize)

	return func(w http.ResponseWriter, r *http.Request) {
		totalReverseRequests.Inc()

		var requestBody reverseRequest
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil || requestBody.Messages == nil {
			error.BadRequestHandler(w, r)
			return
		}
		results := reverseWordsConcurrently(requestBody.Messages, routinesCountSemaphore, cache, cacheKeyMinimumSize)

		response := reverseResponse{Reversed: results}
		json.NewEncoder(w).Encode(response)
	}
}

// reverseWordsConcurrently processes words concurrently while respecting the semaphore that manages the routine limit.
func reverseWordsConcurrently(messages []string, routinesCountSemaphore chan struct{}, cache *Cache, cacheKeyMinimumSize int) []string {
	results := make([]string, len(messages))
	var wg sync.WaitGroup

	for i, word := range messages {
		wg.Add(1)
		routinesCountSemaphore <- struct{}{}
		semaphoreLimitUsage.Inc()

		go func(word string, index int) {
			defer wg.Done()
			defer func() { <-routinesCountSemaphore }()
			defer semaphoreLimitUsage.Dec()

			reversedWord, cacheSizeIncreased := reverseWithCache(word, cache, cacheKeyMinimumSize)
			if cacheSizeIncreased{
				defer messagesCached.Inc()
			}
			results[index] = reversedWord
			messagesReversed.Inc()
		}(word, i)
	}

	wg.Wait()
	return results
}

// reverseWithCache reverses a word, using the cache if available. Returns true if a new word was added to the cache, and the cache size increased.
func reverseWithCache(word string, cache *Cache, cacheKeyMinimumSize int) (string, bool) {
	evictionOccured := false
	cachedValueUsed := true
	wordWasCached := false

	if cached, found := cache.Get(word); found {
		return cached, !evictionOccured && !cachedValueUsed && wordWasCached
	}

	reversed := reverseWord(word)
	cachedValueUsed = false
	if len(word) >= cacheKeyMinimumSize {
		wordWasCached = true
		evictionOccured = cache.Set(word, reversed)
	}
	return reversed, !evictionOccured && !cachedValueUsed && wordWasCached
}

// reverseWord reverses a single word.
func reverseWord(word string) string {
	runes := []rune(word)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
