package msgqueue

import (
	"errors"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jadoint/micro/pkg/logger"
)

// SendMsg send message to Kafka
func SendMsg(topic string, msg []byte) {
	kafkaServer := os.Getenv("KAFKA_SERVER")
	if kafkaServer == "" {
		logger.LogError(errors.New("KAFKA_SERVER is not set"))
		return
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	logger.LogError(err)

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.LogError(fmt.Errorf("Delivery failed: %v", ev.TopicPartition))
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
