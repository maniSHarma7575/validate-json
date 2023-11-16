package jsonparser

import (
	"fmt"
)

func Valid(data []byte) bool {
	content := string(data)

	jsonValidator := NewJSONValidator(content)
	result, err := jsonValidator.ParseJson()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return result
}
