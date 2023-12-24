package main

import (
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	kafkaBroker := "localhost:9092"
	kafkaTopic := "user_rate_limit_exceeded"

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{kafkaBroker}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Println("Shutting down the producer...")
		producer.AsyncClose()
	}()

	for {
		message := &sarama.ProducerMessage{
			Topic: kafkaTopic,
			Value: sarama.StringEncoder("User rate limit exceeded"),
		}
		producer.Input() <- message

		select {
		case success := <-producer.Successes():
			log.Printf("Message %d sent to topic %s, partition %d\n", success.Offset, success.Topic, success.Partition)
		case err := <-producer.Errors():
			log.Println("Failed to send message:", err)
		}
	}
}
