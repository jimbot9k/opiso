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
)

func init() {
	prometheus.MustRegister(semaphoreLimitUsage)
	prometheus.MustRegister(messagesReversed)
	prometheus.MustRegister(totalReverseRequests)
}

func ReverseHandler(routinesCountSemaphore chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		totalReverseRequests.Inc()
		var requestBody reverseRequest
		if err := decoder.Decode(&requestBody); err != nil || requestBody.Messages == nil {
			error.BadRequestHandler(w, r)
			return
		}

		results := make([]string, len(requestBody.Messages))
		var wg sync.WaitGroup

		for i, word := range requestBody.Messages {
			wg.Add(1)
			routinesCountSemaphore <- struct{}{}
			semaphoreLimitUsage.Inc()
			go processWord(word, i, results, &wg, routinesCountSemaphore)
			messagesReversed.Inc()
		}

		wg.Wait()

		result := reverseResponse{Reversed: results}
		json.NewEncoder(w).Encode(result)
	}
}

func processWord(word string, index int, results []string, wg *sync.WaitGroup, routinesCountSemaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-routinesCountSemaphore }()
	defer semaphoreLimitUsage.Dec()

	results[index] = reverseWord(word)
}

func reverseWord(word string) string {
	reversedChars := []rune(word)
	for i, j := 0, len(reversedChars)-1; i < j; i, j = i+1, j-1 {
		reversedChars[i], reversedChars[j] = reversedChars[j], reversedChars[i]
	}
	return string(reversedChars)
}
