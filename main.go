package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AndrewBurian/eventsource"
	redisPubSub "github.com/inklin/golang-sse-test/redis"
)

const (
	serverPort = 3000
)

func main() {

	// start a new stream of events
	stream := eventsource.NewStream()

	// set up Redis Pub
	fmt.Printf("Setting up Redis PubSub\n")
	go redisPubSub.Start(stream)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/event-source", func(rw http.ResponseWriter, req *http.Request) {
		// get topics client wants to subscribe to
		topicsString := req.URL.Query().Get("topics")
		topics := strings.Split(topicsString, ",")

		fmt.Println("creating a new event source client")
		client := eventsource.NewClient(rw, req)
		if client == nil {
			http.Error(rw, "failed to create new client", 400)
			return
		}

		// subscribe client to each topic specified in req query params
		for _, topic := range topics {
			fmt.Println("subscribing client to topic: ", topic)
			stream.Subscribe(topic, client)
		}

		// wait for the client to exit or be shutdown
		client.Wait()
		stream.Remove(client)
		fmt.Println("client removed")
	})

	fmt.Printf("Starting server on port %v\n", serverPort)

	http.ListenAndServe(":"+strconv.Itoa(serverPort), nil)
}
