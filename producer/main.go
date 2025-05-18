package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func tailAndSendToKafka(filePath string, producer *kafka.Producer, topic string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Start from the end of the file
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return fmt.Errorf("error reading file: %w", err)
		}

		// Produce log line to Kafka
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(line),
		}, nil)

		if err != nil {
			fmt.Printf("Failed to send log line to Kafka: %v\n", err)
		} else {
			fmt.Print("Sent to Kafka:", line)
		}
	}
}

func main() {
	filePath := "/var/log/syslog" // Or any other log file path
	topic := "sanaz"

	// Kafka producer config
	config := &kafka.ConfigMap{
		"bootstrap.servers": "192.168.122.10:9092",
	}

	// Create Kafka producer
	producer, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Kafka producer error: %v\n", err)
		os.Exit(1)
	}
	defer producer.Close()

	// Delivery report handler (non-blocking)
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Println("Delivery failed:", ev.TopicPartition.Error)
				}
			}
		}
	}()

	// Tail file and send lines to Kafka
	err = tailAndSendToKafka(filePath, producer, topic)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
