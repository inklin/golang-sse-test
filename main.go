package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AndrewBurian/eventsource"
)

const (
	serverPort = 3000
)

func main() {
	fmt.Printf("Starting server on port %v", serverPort)

	// start a new stream of events
	stream := eventsource.NewStream()

	// create arbitrary topics
	topic := "topic_1"
	topic2 := "topic_2"

	// send event source data events every second
	go func(s *eventsource.Stream) {
		for {
			time.Sleep(time.Second)

			// broadcast tick event to every client
			tickEvent := eventsource.Event{}
			tickMessage := "TICK"
			tickEvent.Data(`{"message":"` + tickMessage + `"}`)
			stream.Broadcast(&tickEvent)

			// publish to topic_1
			// only sent to clients who are subscribed to this topic
			topicEvent := eventsource.Event{}
			topicEvent.Data(`{"message":"Hello from Topic 1","timestamp":"` + time.Now().String() + `"}`)
			stream.Publish(topic, &topicEvent)

			// publish to topic_2
			// only sent to clients who are subscribed to this topic
			topicTwoEvent := eventsource.Event{}
			topicTwoEvent.Data(`{"message":"Hello from Topic 2","timestamp":"` + time.Now().String() + `"}`)
			stream.Publish(topic2, &topicTwoEvent)
		}
	}(stream)

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

	http.ListenAndServe(":"+strconv.Itoa(serverPort), nil)
}
