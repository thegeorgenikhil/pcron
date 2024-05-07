<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![Version](https://img.shields.io/badge/goversion-1.22.x-blue.svg)](https://golang.org)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/thegeorgenikhil/pcron/master/LICENSE.md)

# PCron

PCron is a Go package designed to simplify the integration of cron scheduling with Kafka message production. It provides a structured approach for running scheduled tasks and publishing messages to Kafka topics based on predefined schedules.

## Features

- **Integration with Kafka:** PCron seamlessly integrates with Kafka through the use of the `sarama` package, allowing easy production of messages to Kafka topics.
  
- **Error Handling:** PCron provides an error channel (`errChan`) to capture and handle errors that occur during job execution. This allows for robust error logging or alerting mechanisms to be implemented.

- **Flexible Scheduling:** With PCron, you can define custom schedules using cron expressions, providing flexibility in scheduling tasks according to specific time patterns.

## Installation
Install it in the usual way:

~~~
go get -u github.com/thegeorgenikhil/pcron
~~~

## Usage

Below is a basic example demonstrating how to use PCron to schedule a job that produces messages to a Kafka topic:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/thegeorgenikhil/pcron"
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

// This method is called by the producer cron to run the job
// It returns a slice of producer messages to be sent to Kafka
func (hwj *HelloWorldJob) Run() ([]*sarama.ProducerMessage, error) {
	messages := []*sarama.ProducerMessage{
		{
			Topic: TopicName,
			Key:   sarama.StringEncoder("key"),
			Value: sarama.StringEncoder("Hello World!"),
		},
	}

	return messages, nil
}

func main() {
	job := &HelloWorldJob{}

    // Before running the producer cron, make sure you create the topic in Kafka.
    // An example for creating a topic is given in the `examples` folder.
	createTopic()

	// Create producer cron config
	config := pcron.NewConfig(Name, Schedule, job, BrokerURLs)

	// Create producer cron
	producerCron, err := pcron.New(config)
	if err != nil {
		log.Fatalf("Error creating producer cron: %v", err)
	}

	// Start the cron
	err = producerCron.StartCron()
	if err != nil {
		log.Fatalf("Error starting producer cron: %v", err)
	}
	defer producerCron.StopCron()

	// Handle errors from the job in a separate goroutine
	go func() {
		for err := range producerCron.GetErrorChan() {
			log.Printf("[ERROR]: %v", err)
		}
	}()

	select {}
}
```