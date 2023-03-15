package main

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func main() {
	var result Response
	if err := json.Unmarshal([]byte(`{"errors": [{"message": "ERROR at Object.Field..."}]}`), &result); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", result)
}
