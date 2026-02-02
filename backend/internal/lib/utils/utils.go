package utils

import (
	"encoding/json"
	"fmt"
)

// PrintJSON pretty prints the given value as JSON to the console.
func PrintJSON(v interface{}) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println("JSON:", string(json))
}
