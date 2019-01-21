package main

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	serverPort = 3000
)

func feedHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("feed endpoint")
	// add logic for subscribing to multiple sse "rooms" here
}

func main() {
	fmt.Printf("Starting server on port %v", serverPort)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/feed", feedHandler)

	http.ListenAndServe(":"+strconv.Itoa(serverPort), nil)
}
