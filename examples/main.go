package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/thegeorgenikhil/cronp"
)

const (
	Name = "hello-world-cron"
	Schedule = "* * * * *" // Run every minute
	TopicName = "hello-world"
)

var (
	BrokerURLs = []string{"localhost:9092"}
)

type HelloWorldJob struct{}

func (hwj *HelloWorldJob) Run() ([]*sarama.ProducerMessage, error) {
	messages := []*sarama.ProducerMessage{
		{
			Topic: TopicName,
			Key:   sarama.StringEncoder("key"),
			Value: sarama.StringEncoder("Hello World!"),
		},
	}

	// Randomly generate error if minute is divisible by 4
	// This is to show error handling in action
	if time.Now().Minute()%4 == 0 {
		return nil, fmt.Errorf("ouch! error occurred")
	}
	return messages, nil
}

func main() {
	job := &HelloWorldJob{}

	createTopic()

	// Create cron producer config
	config := cronp.NewConfig(Name, Schedule, job, BrokerURLs)

	// Create cron producer
	cronProducer, err := cronp.New(config)
	if err != nil {
		log.Fatalf("Error creating cron producer: %v", err)
	}

	// Start the cron
	err = cronProducer.StartCron()
	if err != nil {
		log.Fatalf("Error starting cron producer: %v", err)
	}
	defer cronProducer.StopCron()

	// Handle errors from the job in a separate goroutine
	go func() {
		for err := range cronProducer.GetErrorChan() {
			log.Printf("[ERROR]: %v", err)
		}
	}()

	select {}
}

func createTopic() {
	admin, err := sarama.NewClusterAdmin(BrokerURLs, sarama.NewConfig())
	if err != nil {
		log.Fatalf("Error creating cluster admin: %v", err)
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing topics: %v", err)
	}

	if _, ok := topics[TopicName]; !ok {
		err = admin.CreateTopic(TopicName, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			log.Printf("Error creating topic: %v", err)
		}
		log.Printf("Topic %s created", TopicName)
	} else {
		log.Printf("Topic %s already exists", TopicName)
	}	
}
