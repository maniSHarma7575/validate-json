package jsonparser

import (
	"os"
	"strconv"
	"testing"
	"validatejson/pkg/jsonparser"

	"github.com/stretchr/testify/assert"
)

func content(fileNumber int, fileType string) string {
	content, _ := os.ReadFile("../data/" + fileType + strconv.Itoa(fileNumber) + ".json")
	return string(content)
}
func TestParseJson(t *testing.T) {
	jsonValidator := jsonparser.NewJSONValidator(content(3, "pass"))
	result, err := jsonValidator.ParseJson()
	assert.Equal(t, true, result)
	assert.Nil(t, err)
}

func TestParseJsonForInvalidJson(t *testing.T) {
	var tests = []struct {
		name         string
		fileNumber   int
		errorMessage string
	}{
		{"json payload should be object or array", 1, "a json payload should be an object or array"},
		{"Unclosed array", 2, "closing bracket for the array is missing"},
		{"keys must be quoted", 3, "object key name is not a string"},
		{"array given value is not correct", 4, "array value is not correct"},
		{"Comma after the close", 7, "invalid json: extra character after closing bracket"},
		{"Object key should be string", 9, "object key name is not a string"},
		{"Illegal expression passed", 11, "invalid json: Expected ','"},
		{"Object value should be valid", 12, "object value is not in correct format"},
		{"Number should be in correct format", 13, "number cannot have leading zeros"},
		{"Illegal backslash escape sequence", 15, "invalid json: Illegal backslash escape sequence"},
		{"Illegal string", 25, "invalid json: Illegal character tab character or new line character"},
		{"Not a valid number", 29, "not a valid number"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			jsonValidator := jsonparser.NewJSONValidator(content(test.fileNumber, "fail"))
			result, err := jsonValidator.ParseJson()
			assert.Equal(t, false, result)
			assert.NotNil(t, err)
			assert.Equal(t, test.errorMessage, err.Error())
		})
	}
}
