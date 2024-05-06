package cronp

const (
	// __defaultTopicPartitions is the default number of partitions for a topic
	__defaultTopicPartitions = 1

	// __defaultTopicReplicationFactor is the default replication factor for a topic
	__defaultTopicReplicationFactor = 1
)

type config struct {
	name string
	schedule string
	job Job
	brokerURLs []string
	topicName string
	topicPartitions int
	topicReplicationFactor int
}

func NewConfig(name string, schedule string, job Job, brokerURL []string, topic string) *config {
	return &config{
		name: name,
		schedule: schedule,
		job: job,
		brokerURLs: brokerURL,
		topicName: topic,
		topicPartitions: __defaultTopicPartitions,
		topicReplicationFactor: __defaultTopicReplicationFactor,
	}
}

func (c *config) SetTopicPartitions(partitions int) {
	c.topicPartitions = partitions
}

func (c *config) SetTopicReplicationFactor(replicationFactor int) {
	c.topicReplicationFactor = replicationFactor
}