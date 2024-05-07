package pcron

import (
	"github.com/IBM/sarama"
	"github.com/robfig/cron/v3"
)

// ProducerCron contains the configuration and associated cron meth
type ProducerCron struct {
	// config holds the configuration for the cron job.
	config   *Config
	// cron is the cron scheduler.
	cron     *cron.Cron
	// producer is the Kafka producer which publishes messages.
	producer sarama.AsyncProducer
	// errChan is the error channel which receives errors from the job, that can be used to log errors or send alerts.
	errChan  chan error
}

// NewProducerCron creates a new ProducerCron instance.
func New(cfg *Config) (*ProducerCron, error) {
	producer, err := sarama.NewAsyncProducer(cfg.BrokerURLs, cfg.ProducerConfig)
	if err != nil {
		return nil, err
	}

	return &ProducerCron{
		config:   cfg,
		producer: producer,
		errChan:  make(chan error),
	}, nil
}

// StartCron starts the cron scheduler and runs the job at the specified schedule.
func (cp *ProducerCron) StartCron() error {
	cp.cron = cron.New()
	_, err := cp.cron.AddFunc(cp.config.Schedule, func() {
		messages, err := cp.config.Job.Run()
		if err != nil {
			// Send error to error channel
			// This can be used to log errors or send alerts
			cp.errChan <- err
			return
		}

		for _, message := range messages {
			cp.producer.Input() <- message
		}
	})
	if err != nil {
		return err
	}

	cp.cron.Start()
	return nil
}

// GetErrorChan returns the error channel.
func (cp *ProducerCron) GetErrorChan() <-chan error {
	return cp.errChan
}

// StopCron stops the cron scheduler.
func (cp *ProducerCron) StopCron() {
	cp.producer.Close()
	if cp.cron != nil {
		cp.cron.Stop()
	}
}
