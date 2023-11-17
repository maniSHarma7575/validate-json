package jsonparser

import (
	"os"
	"strconv"
	"testing"
	"validatejson/pkg/jsonparser"

	"github.com/stretchr/testify/assert"
)

func TestJSONValid(t *testing.T) {
	for i := range [3]int{} {
		content, _ := os.ReadFile("../data/pass" + strconv.Itoa(i+1) + ".json")
		result := jsonparser.Valid([]byte(content))
		assert.Equal(t, true, result)
	}
}

func TestJSONInValid(t *testing.T) {
	for i := range [31]int{} {
		if i != 17 {
			content, _ := os.ReadFile("../data/fail" + strconv.Itoa(i+1) + ".json")
			result := jsonparser.Valid([]byte(content))
			assert.Equal(t, false, result)
		}
	}
}

func BenchmarkJsonValid(b *testing.B) {
	content, _ := os.ReadFile("../data/pass1.json")
	for i := 0; i < b.N; i++ {
		jsonparser.Valid([]byte(content))
	}
}
