package main

import (
	"fmt"

	sse "github.com/r3labs/sse/v2"
)

func main() {
	client := sse.NewClient("http://localhost:8844/stream?")
	client.Subscribe("message", func(msg *sse.Event) {
		// Got some data!
		fmt.Println(string(msg.Data))
	})
}
