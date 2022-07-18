# Cloud Emulator Init

Initialize resources in cloud emulators

---

## Help Text

```
cloud-emulator-init --help
initialize cloud emulation resources

Usage:
  cloud-emulator-init [flags]

Flags:
      --config-file string   file containing resources to create (default "conf.yaml")
  -h, --help                 help for cloud-emulator-init
```

## Running locally

### Install binary

```shell
go install github.com/jimmykodes/cloud-emulator-init@latest

cloud-emulator-init
```

### Docker

```shell
docker pull jimmykodes/cloud-emulator-init:latest

docker run \
  -v ${PWD}/conf.yaml:/etc/cloud-emulator-init/conf.yaml \
  -e CONFIG_FILE=/etc/cloud-emulator-init/conf.yaml \
  jimmykodes/cloud-emulator-init
```

## Supported Emulators

- Localstack (AWS)
- Google cloud-sdk emulators (GCP)

## Supported Objects

### AWS

- SQS Queues

### GCP

- Pubsub topics/subscriptions

## Config

Resources are defined in a yaml file (default `conf.yaml`).

#### Example

```yaml
aws:
  emulatorUrl: "http://localhost:4566"
  sqs:
    - my-queue
    - my-other-queue
gcp:
  project: my-project
  pubsub:
    - name: topic1
      subscriptions:
        - topic1sub1
        - topic1sub2
    - name: topic2
      subscriptions:
        - topic2sub1
```

## Example Usage

### Docker Compose

```yaml
version: '3'
services:
  sqs:
    image: localstack/localstack
    ports:
      - "4566:4566"
  pubsub:
    image: google/cloud-sdk
    command: "gcloud beta emulators pubsub start --host-port 0.0.0.0:8085 --log-http --verbosity debug --user-output-enabled"
    ports:
      - "8085:8085"
  init:
    image: jimmykodes/cloud-emulator-init:latest
    volumes:
      - ./conf.yaml:/etc/cloud-emulator-init/conf.yaml
    environment:
      PUBSUB_EMULATOR_HOST: pubsub:8085
      CONFIG_FILE: /etc/cloud-emulator-init/conf.yaml
    restart: on-failure
    depends_on:
      - sqs
      - pubsub
```

