package main

import (
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

var rabbitConn *amqp.Connection

func main() {
	setupRabbitMQProducer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order Service"))
	})
	http.HandleFunc("/placeOrder", placeOrderHandler)

	log.Fatal(http.ListenAndServe(":5000", nil))
}

func setupRabbitMQProducer() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	log.Println("Connected to RabbitMQ successfully")
}

func placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publishOrderToRabbitMQ(order)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order placed successfully"))
}

func publishOrderToRabbitMQ(order Order) {
	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"processed_orders", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return
	}

	body, err := json.Marshal(order)
	if err != nil {
		log.Println("Failed to marshal order:", err)
		return
	}

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key (queue name)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("Failed to publish a message:", err)
		return
	}
	log.Println("Order published to RabbitMQ successfully")
}
