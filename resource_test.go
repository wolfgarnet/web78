package core

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestResourceStruct(t *testing.T) {
	jsonData := []byte(`{"id": "user-1", "name":"wolle"}`)
	fmt.Println("JSON", string(jsonData))


	println(jsonData)
	var res BaseResource
	if err := json.Unmarshal(jsonData, &res); err != nil {
		println("ERROR", err.Error())
	}

	fmt.Printf("RES___----> %v\n", res.Metadata.Title)
}