# pubsub-devsub

## Installation

To install cli, simply run:
```
go get github.com/groovenauts/pubsub-devsub
```

Make sure your PATH includes the $GOPATH/bin directory so your commands can be easily used:

```
export PATH=$PATH:$GOPATH/bin
```

# pubsub-devsub

## Installation

Download the file from https://github.com/groovenauts/pubsub-devsub/releases and put it somewhere on PATH.

## Usage

```bash
$ pubsub-devsub
NAME:
   pubsub-devsub - github.com/groovenauts/pubsub-devsub

USAGE:
   pubsub-devsub [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --project value       GCS Project ID [$GCP_PROJECT, $PROJECT]
   --subscription value  Subscription [$SUBSCRIPTION]
   --interval value      Interval to pull (default: 10)
   --help, -h            show help
   --version, -v         print the version
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
