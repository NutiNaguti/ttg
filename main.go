package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	client := NewClient()
	update, _ := client.Receive()
	var data = make(map[string]interface{})
	json.Unmarshal(update, &data)
	fmt.Println(data)
}
