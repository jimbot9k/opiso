package reverse

import (
	"encoding/json"
	"gilmour/opiso/internal/error"
	"net/http"
)

type ReversedWord struct {
	Original string `json:"original"`
	Reversed string `json:"reversed"`
}

func ReverseHandler(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")
	if len(word) == 0 {
		error.BadRequestHandler(w, r)
		return
	}
	result := ReversedWord{Original: word, Reversed: reverseWord(word)}

	json.NewEncoder(w).Encode(result)
}

func reverseWord(word string) string {
	reversedChars := []rune(word)
	for i, j := 0, len(reversedChars)-1; i < j; i, j = i+1, j-1 {
		reversedChars[i], reversedChars[j] = reversedChars[j], reversedChars[i]
	}
	return string(reversedChars)
}
