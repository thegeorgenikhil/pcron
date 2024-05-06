package cronp

import (
	"context"
	"net"
	"strconv"

	"github.com/robfig/cron/v3"
	kafka "github.com/segmentio/kafka-go"
)

type Job interface {
	Run() ([]kafka.Message, error)
}

type cronProducer struct {
	config *config
	cron  *cron.Cron
	writer *kafka.Writer
	errCh  chan error
}

func NewCronProducer(ctx context.Context, cfg *config) (*cronProducer, error) {
	err := createTopics(
		context.Background(),
		cfg.brokerURLs[0],
		cfg.topicName,
		cfg.topicPartitions,
		cfg.topicReplicationFactor,
	)
	if err != nil {
		return nil, err
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.brokerURLs,
		Topic:    cfg.topicName,
		RequiredAcks: 1,
	})

	return &cronProducer{
		config: cfg,
		writer: writer,
	}, nil
}

func (cp *cronProducer) StartCron() error {
	cp.cron = cron.New(cron.WithSeconds())
	_, err := cp.cron.AddFunc(cp.config.schedule, func() {
		messages, err := cp.config.job.Run()
		if err != nil {
			cp.errCh <- err
		}

		err = cp.writer.WriteMessages(context.Background(), messages...)
		if err != nil {
			cp.errCh <- err
		}
	})
	if err != nil {
		return err
	}

	cp.cron.Start()
	return nil
}

func (cp *cronProducer) StopCron() {
	cp.writer.Close()
	cp.cron.Stop()
}

func (cp *cronProducer) GetErrors() <-chan error {
	return cp.errCh
}

func createTopics(ctx context.Context, brokerUrl string, topicName string, partitions, replicationFactor int) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokerUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, _ := conn.Controller()

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		return err
	}

	return nil
}
