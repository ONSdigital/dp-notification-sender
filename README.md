dp-dataset-exporter
================

⚠️ 10% PROJECT - DO NOT DEPLOY

Takes a completed filter job message and produces notification

### Getting started

* Run `make debug`


### Configuration

An overview of the configuration options available, either as a table of
environment variables, or with a link to a configuration guide.

| Environment variable        | Default                              | Description
| --------------------------- | ------------------------------------ | -----------
| BIND_ADDR                   | :23456                               | The host and port to bind to
| KAFKA_ADDR                  | localhost:9092                       | The address of Kafka
| FILTER_JOB_CONSUMER_TOPIC   | filter-job-submitted                 | The name of the topic to consume messages from
| FILTER_JOB_CONSUMER_GROUP   | dp-dataset-exporter                  | The consumer group this application to consume filter job messages
| HEALTHCHECK_INTERVAL        | time.Minute                          | How often to run a health check
| GRACEFUL_SHUTDOWN_TIMEOUT   | time.Second * 10                     | How long to wait for the service to shutdown gracefully


### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
