package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var msgChan chan string

func main() {
	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))

	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.HandleFunc("/event", sseHandler)
	router.HandleFunc("/ping", ping)

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// ping will send the current time to the client
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sending time to client")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if msgChan != nil {
		msg := time.Now().Format("15:04:05") + "ping from " + r.RemoteAddr
		msgChan <- msg
	}

	w.Write([]byte(time.Now().Format("15:04:05") + " ping sent"))
}

// sseHandler will send events to the client
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	msgChan = make(chan string)

	defer func() {
		close(msgChan)
		msgChan = nil
		fmt.Println("Client closed connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not initialize flusher")
	}

	for {
		select {
		case msg := <-msgChan:
			fmt.Println("Sending message to client: ", msg)
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("done")
			return
		}
	}
}
