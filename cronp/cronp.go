package cronp

import (
	"github.com/IBM/sarama"
	"github.com/robfig/cron/v3"
)

// CronProducer contains the configuration and associated cron meth
type CronProducer struct {
	config   *Config
	cron     *cron.Cron
	producer sarama.AsyncProducer
	errChan  chan error
}

// NewCronProducer creates a new CronProducer instance.
func NewCronProducer(cfg *Config) (*CronProducer, error) {
	producer, err := sarama.NewAsyncProducer(cfg.BrokerURLs, cfg.ProducerConfig)
	if err != nil {
		return nil, err
	}

	return &CronProducer{
		config:   cfg,
		producer: producer,
		errChan:  make(chan error),
	}, nil
}

// StartCron starts the cron scheduler and runs the job at the specified schedule.
func (cp *CronProducer) StartCron() error {
	cp.cron = cron.New()
	_, err := cp.cron.AddFunc(cp.config.Schedule, func() {
		messages, err := cp.config.Job.Run()
		if err != nil {
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
func (cp *CronProducer) GetErrorChan() <-chan error {
	return cp.errChan
}

// StopCron stops the cron scheduler.
func (cp *CronProducer) StopCron() {
	cp.producer.Close()
	if cp.cron != nil {
		cp.cron.Stop()
	}
}
