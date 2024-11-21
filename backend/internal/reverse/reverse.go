package reverse

import (
	"encoding/json"
	"github.com/jimbot9k/opiso/internal/error"
	"net/http"
	"sync"
)

type reverseRequest struct {
	Messages []string `json:"messages"`
}

type reverseResponse struct {
	Reversed []string `json:"reversed"`
}

func ReverseHandler(routinesCountSemaphore chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
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
			go processWord(word, i, results, &wg, routinesCountSemaphore)
		}

		wg.Wait()

		result := reverseResponse{Reversed: results}
		json.NewEncoder(w).Encode(result)
	}
}

func processWord(word string, index int, results []string, wg *sync.WaitGroup, routinesCountSemaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-routinesCountSemaphore }()

	results[index] = reverseWord(word)
}

func reverseWord(word string) string {
	reversedChars := []rune(word)
	for i, j := 0, len(reversedChars)-1; i < j; i, j = i+1, j-1 {
		reversedChars[i], reversedChars[j] = reversedChars[j], reversedChars[i]
	}
	return string(reversedChars)
}
