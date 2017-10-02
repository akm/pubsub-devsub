# pubsub-devsub

## Installation

Download the file from https://github.com/groovenauts/pubsub-devsub/releases and put it somewhere on PATH.

## Usage

### Subscribe with ACK

```
$ pubsub-devsub subscribe -h
NAME:
   pubsub-devsub subscribe - Subscribe and show message with ACK

USAGE:
   pubsub-devsub subscribe [command options] [arguments...]

OPTIONS:
   --interval value                Interval to pull (default: 10)
   --max-messages value, -m value  Max messages per pull (default: 10)
   --project value                 GCS Project ID [$GCP_PROJECT, $PROJECT]
   --return-immediately, -r        Return result immediately on pull
   --verbose, -V                   Show debug logs
```

For example

```bash
$ pubsub-devsub subscribe -V -m 1 projects/project1/subscriptions/sub1
```


You can stop this by `Ctrl+C` or `-INT` signal.

### Inspect messages

```
$ pubsub-devsub inspect -h
NAME:
   pubsub-devsub inspect - Pull and show messages without ACK

USAGE:
   pubsub-devsub inspect [command options] [arguments...]

OPTIONS:
   --interval value                Interval to pull (default: 10)
   --max-messages value, -m value  Max messages per pull (default: 10)
   --project value                 GCS Project ID [$GCP_PROJECT, $PROJECT]
   --return-immediately, -r        Return result immediately on pull
   --verbose, -V                   Show debug logs
```
For example

```bash
$ pubsub-devsub subscribe -V -r projects/project1/subscriptions/sub1
```


### Primitive Usage

```bash
$ pubsub-devsub
NAME:
   pubsub-devsub - development tool for Google Cloud Pubsub

USAGE:
   pubsub-devsub [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     inspect, i    Pull and show messages without ACK
     subscribe, s  Subscribe and show message with ACK
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ack, -A                       Send ACK for received message
   --follow, -f                    Keep subscribing
   --interval value                Interval to pull (default: 10)
   --max-messages value, -m value  Max messages per pull (default: 10)
   --project value                 GCS Project ID [$GCP_PROJECT, $PROJECT]
   --return-immediately, -r        Return result immediately on pull
   --verbose, -V                   Show debug logs
   --help, -h                      show help
   --version, -v                   print the version
```

## How to biuld

```
$ make setup
$ make build
```

## Release

```
$ make release
```
