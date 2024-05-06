package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	cronp "github.com/thgeorgenikhil/cronp/cronp"
)

type wsjJob struct {
}

func (w *wsjJob) Run() ([]kafka.Message, error) {
	return []kafka.Message{
		{
			Key:   []byte("wsj"),
			Value: []byte("wsj"),
		},
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The following code snippet demonstrates how to use the cronp package to create a cron job that runs every 10 seconds.
	cfg := cronp.NewConfig("wsj-cron", "*/10 * * * * *", &wsjJob{}, []string{"localhost:9092"}, "wsj")
	cfg.SetTopicPartitions(10)
	cfg.SetTopicPartitions(3)

	cronProducer, err := cronp.NewCronProducer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = cronProducer.StartCron()
	if err != nil {
		log.Fatal(err)
	}
	defer cronProducer.StopCron()
	for {}
}
