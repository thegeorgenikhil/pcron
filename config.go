package cronp

import "github.com/IBM/sarama"

// Config holds the configuration for a cron job.
type Config struct {
	// Name is the name of the cron job.
	Name string
	// Schedule is the cron schedule for the job.
	Schedule string
	// Job is the job to be executed. It has a Run method where the actual job logic is implemented and returns a list of messages to be published.
	Job Job
	// BrokerURLs is the list of Kafka broker URLs.
	BrokerURLs []string
	// ProducerConfig is the configuration for the Kafka producer.
	ProducerConfig *sarama.Config
}

// NewConfig creates a new Config instance.
func NewConfig(name, schedule string, job Job, brokerURls []string, producerConfig ...*sarama.Config) *Config {
	cfg := &Config{
		Name:       name,
		Schedule:   schedule,
		Job:        job,
		BrokerURLs: brokerURls,
	}

	producerCfg := sarama.NewConfig()
	if len(producerConfig) > 0 {
		producerCfg = producerConfig[0]
	}

	cfg.ProducerConfig = producerCfg

	return cfg
}
