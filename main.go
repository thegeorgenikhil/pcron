package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	cronp "github.com/thgeorgenikhil/cronp/cronp"
)

type wsjJob struct {
}

func (w *wsjJob) Run() ([]*sarama.ProducerMessage, error) {
	messages := []*sarama.ProducerMessage{
		{
			Topic: "wsj",
			Key:   sarama.StringEncoder("key"),
			Value: sarama.StringEncoder("Hello World!"),
		},
	}

	// Randomly generate error if minute is even
	if time.Now().Minute()%2 == 0 {
		return nil, fmt.Errorf("error")
	}
	return messages, nil
}

func main() {
	job := &wsjJob{}
	config := cronp.NewConfig("wsj", "* * * * *", job, []string{"localhost:9092"}, "wsj")
	cronProducer, err := cronp.NewCronProducer(config)
	if err != nil {
		log.Fatalf("error creating cron producer: %v", err)
	}

	err = cronProducer.StartCron()
	if err != nil {
		log.Fatalf("error starting cron producer: %v", err)
	}

	go func() {
		for err := range cronProducer.GetErrorChan() {
			log.Printf("error: %v", err)
		}
	}()

	select {}
}
