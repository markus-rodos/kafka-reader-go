# Event processing service

### Setup
To set up golang dependencies: 
`go mod tidy`

### Local run
To run the application locally, one needs to do the following:

`docker compose up` - starts kafka as docker container

`go run main.go` - starts an application

`$KAFKA_HOME/bin/kafka-console-producer.sh --topic events.all --bootstrap-server localhost:9092 < events.json` - produces events from events.json
where, KAFKA_HOME is an env variable pointing to kafka installation directory; events.json file with events to put to the topic


