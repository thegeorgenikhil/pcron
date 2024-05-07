<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![Version](https://img.shields.io/badge/goversion-1.22.x-blue.svg)](https://golang.org)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/thegeorgenikhil/cronp/master/LICENSE.md)

# CronP

CronP is a Go package designed to simplify the integration of cron scheduling with Kafka message production. It provides a structured approach for running scheduled tasks and publishing messages to Kafka topics based on predefined schedules.

## Features

- **Integration with Kafka:** CronP seamlessly integrates with Kafka through the use of the `sarama` package, allowing easy production of messages to Kafka topics.
  
- **Error Handling:** CronP provides an error channel (`errChan`) to capture and handle errors that occur during job execution. This allows for robust error logging or alerting mechanisms to be implemented.

- **Flexible Scheduling:** With CronP, you can define custom schedules using cron expressions, providing flexibility in scheduling tasks according to specific time patterns.

## Installation
Install it in the usual way:

~~~
go get -u github.com/thegeorgenikhil/cronp
~~~

## Usage

Below is a basic example demonstrating how to use CronP to schedule a job that produces messages to a Kafka topic:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/thgeorgenikhil/cronp"
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

// This method is called by the cron producer to run the job
// It returns a slice of producer messages to be sent to Kafka
func (hwj *HelloWorldJob) Run() ([]*sarama.ProducerMessage, error) {
	messages := []*sarama.ProducerMessage{
		{
			Topic: TopicName,
			Key:   sarama.StringEncoder("key"),
			Value: sarama.StringEncoder("Hello World!"),
		},
	}

	// Randomly generate error if minute is divisible by 3
	// This is to show error handling in action
	if time.Now().Minute()%4 == 0 {
		return nil, fmt.Errorf("ouch! error occurred")
	}
	return messages, nil
}

func main() {
	job := &HelloWorldJob{}

    // Before running the cron producer, make sure you create the topic in Kafka. An example for creating a topic is given in the `examples` folder.
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
```