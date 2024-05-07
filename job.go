package pcron

import "github.com/IBM/sarama"

// Job defines the interface for jobs to be executed.
type Job interface {
	Run() ([]*sarama.ProducerMessage, error)
}
