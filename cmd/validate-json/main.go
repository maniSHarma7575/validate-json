package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"validatejson/pkg/jsonparser"
)

func main() {
	fileFlag := flag.String("file", "", "File name")

	flag.Parse()

	fileName := strings.TrimSpace(*fileFlag)

	content, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", fileName, err)
		os.Exit(0)
	}

	fmt.Printf("Valid: %t\n", jsonparser.Valid([]byte(content)))
}
