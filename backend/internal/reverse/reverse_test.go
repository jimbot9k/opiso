package reverse

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReverseWord(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"", ""},
		{"a", "a"},
		{"racecar", "racecar"},
		{"I am the dog, the big bad doggggggggggg!", "!gggggggggggod dab gib eht ,god eht ma I"},
	}

	for _, test := range tests {
		output := reverseWord(test.input)
		if output != test.expected {
			t.Errorf("Expected %q, got %q", test.expected, output)
		}
	}
}

func TestReverseWithCache(t *testing.T) {
	cache := NewCache(2)

	tests := []struct {
		input          string
		expectedOutput string
		expectedCache  bool
	}{
		{"hello", "olleh", true},  // New entry
		{"world", "dlrow", true},  // New entry
		{"hello", "olleh", false}, // Cached entry
	}

	for _, test := range tests {
		output, cacheIncreased := reverseWithCache(test.input, cache, 3)
		if output != test.expectedOutput {
			t.Errorf("Expected output %q, got %q", test.expectedOutput, output)
		}
		if cacheIncreased != test.expectedCache {
			t.Errorf("Expected cache increased %v, got %v", test.expectedCache, cacheIncreased)
		}
	}
}

func TestReverseWordsConcurrently(t *testing.T) {
	routinesCountSemaphore := make(chan struct{}, 2) // Limit to 2 concurrent routines
	cache := NewCache(10)

	input := []string{"hello", "world", "golang", "rocks"}
	expected := []string{"olleh", "dlrow", "gnalog", "skcor"}

	output := reverseWordsConcurrently(input, routinesCountSemaphore, cache, 3)

	for i, result := range output {
		if result != expected[i] {
			t.Errorf("At index %d, expected %q, got %q", i, expected[i], result)
		}
	}
}

func TestReverseHandlerWithCache(t *testing.T) {
	routinesCountSemaphore := make(chan struct{}, 2)
	cacheSize := 10
	cacheKeyMinimumSize := 3

	handler := ReverseHandlerWithCache(routinesCountSemaphore, cacheSize, cacheKeyMinimumSize)

	t.Run("Valid request", func(t *testing.T) {
		requestBody := reverseRequest{
			Messages: []string{"My", "name", "is", "Jonas", "I'm", "carrying", "the", "wheel", "Thanks", "for", "all", "you've", "shown", "us", "This", "is", "how", "we", "feel"},
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/reverse", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		var response reverseResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		expected := []string{"yM", "eman", "si", "sanoJ", "m'I", "gniyrrac", "eht", "leehw", "sknahT", "rof", "lla", "ev'uoy", "nwohs", "su", "sihT", "si", "woh", "ew", "leef"}
		for i, result := range response.Reversed {
			if result != expected[i] {
				t.Errorf("At index %d, expected %q, got %q", i, expected[i], result)
			}
		}
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/reverse", bytes.NewReader([]byte("this-is-clearly-not-json")))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Invalid request body but still json", func(t *testing.T) {

		requestBody := struct {
			BadKey string `json:"badKey"`
		}{BadKey: "bad value"}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/reverse", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, gota %d", rec.Code)
		}
	})

	t.Run("Invalid request body invalid messages", func(t *testing.T) {

		requestBody := struct {
			Messages string `json:"messages"`
		}{Messages: "not an array"}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/reverse", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, gota %d", rec.Code)
		}
	})
}
