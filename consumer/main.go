package main

import (
	"fmt"
	"regexp"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	re := regexp.MustCompile(`:+\s.*`)

	// Set up Kafka consumer configuration
	fmt.Println("start listening app1")
	config := &kafka.ConfigMap{

		"bootstrap.servers": "192.168.122.10:9092", // broker addr
		"group.id":          "console-consumer-3",  // consumer group
		"auto.offset.reset": "beginning",           // earliest offset

	}

	// Create Kafka consumer

	consumer, err := kafka.NewConsumer(config)

	if err != nil {
		panic(err)
	}

	// subscribe to target topic
	consumer.SubscribeTopics([]string{"sanaz"}, nil)
	// Consume messages
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			mmd := re.FindString(string(msg.Value))
			fmt.Printf("Received message: %s\n", mmd)
			// fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Error while consuming message: %v (%v)\n", err, msg)
		}
	}
	// Close Kafka consumer
}
