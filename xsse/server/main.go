package main

import (
	"fmt"
	"net/http"
	"time"

	sse "github.com/r3labs/sse/v2"
)

func main() {
	server := sse.New()
	defer server.Close()

	// 定义一个 channel
	var closeCh chan struct{}

	server.AutoReplay = false
	server.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}
	server.CreateStream("message")

	// Create a new Mux and set the handler
	mux := http.NewServeMux()

	// 配置静态文件服务
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	mux.HandleFunc("/stream", server.ServeHTTP)

	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if closeCh == nil {
			closeCh = make(chan struct{})
			go publish(server, closeCh)
			return
		}

		w.Write([]byte("Already started!"))
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if closeCh != nil {
			close(closeCh)
			closeCh = nil
		}

		w.Write([]byte("Please start first!"))
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	fmt.Println("Server running on :8844")

	http.ListenAndServe(":8844", mux)
}

func publish(server *sse.Server, closeCh chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("publishing")
			server.Publish("message", &sse.Event{
				Event: []byte("message"),
				Data:  []byte(fmt.Sprintf("%s", time.Now().String())),
			})
		case <-closeCh:
			server.Publish("message", &sse.Event{
				Event: []byte("close"),
				Data:  []byte("stop..."),
			})
			return
		}
	}
}
