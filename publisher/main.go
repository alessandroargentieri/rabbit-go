package main

import (
	"fmt"
	"log"
	"os"
    "net/http"
    "encoding/json"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

var publishedMsgs []string

func startPublisher() {

	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	url := os.Getenv("RABBITMQ_URL")
	
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", user, password, url))
	if err != nil {
		log.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"events", // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Println(err)
	}

	for i := 1; i < 30; i++ {
		time.Sleep(5 * time.Second)

		message := fmt.Sprintf("Message n. %d: %s", i, randomString(6))
		// attempt to publish a message to the queue!
		err = ch.Publish(
			"events",
			"",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		if err != nil {
			log.Println(err)
		}
		publishedMsgs = append(publishedMsgs, message)
		log.Println("Successfully Published Message to Queue: ", message)
	}
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func main() {
	// start RabbitMQ publisher on an everlasting parallel goroutine
	go startPublisher()

	// start HTTP server
	http.HandleFunc("/publisher", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(Statistics{len(publishedMsgs), getLatestMessages()})
		log.Println(string(resp))
        w.Write([]byte(resp))
	})            
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

type Statistics struct {
	Sent int `json:"sent"`
	Latests  []string `json:"latests"`
}

func getLatestMessages() []string {
	if len(publishedMsgs) <= 5 {
		return publishedMsgs
	} 
    return publishedMsgs[5:]
}