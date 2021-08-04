package main

import (
	"fmt"
	"log"
	"os"
    "net/http"
    "encoding/json"
	"github.com/streadway/amqp"
)

var consumedMsgs []string

func startConsumer() {

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
		fmt.Println(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"events", // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			msg := fmt.Sprintf("Received Message: %s", d.Body)
			log.Println(msg)
			consumedMsgs = append(consumedMsgs, msg)
		}
	}()

	log.Println("Successfully Connected to our RabbitMQ Instance")
	log.Println(" [*] - Waiting for messages")
	<-forever
}

func main() {
	// start RabbitMQ listener and consumer on an everlasting parallel goroutine
	go startConsumer()

	// start HTTP server
	http.HandleFunc("/consumer", func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(Statistics{len(consumedMsgs), getLatestMessages()})
		log.Println(string(resp))
        w.Write([]byte(resp))
	})            
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

type Statistics struct {
	Received int `json:"received"`
	Latests  []string `json:"latests"`
}

func getLatestMessages() []string {
	if len(consumedMsgs) <= 5 {
		return consumedMsgs
	} 
    return consumedMsgs[5:]
}