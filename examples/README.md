# PCron - Example

This example demonstrates how to use PCron to schedule a job, which runs every minute and produces messages to a Kafka topic called `hello-world`. The code also creates a topic at the start to ensure that the topic exists before t2he producer starts producing messages.

Error handling is also demonstrated in the example. The job randomly generates an error if the current minute is divisible by 4.

## How to Run?

1. Run Kafka and Zookeeper using the `docker-compose.yml` file given

```bash
docker-compose up -d
```

2. In a separate terminal, run the `main.go` file

```bash
go run main.go
```