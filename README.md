# pubsub-devsub

## Installation

To install cli, simply run:
```
go get github.com/akm/pubsub-devsub
```

Make sure your PATH includes the $GOPATH/bin directory so your commands can be easily used:

```
export PATH=$PATH:$GOPATH/bin
```

## Usage

```bash
$ pubsub-devsub
NAME:
   pubsub-devsub - github.com/akm/pubsub-devsub

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

At first, install gom.
```
go get github.com/mattn/gom
```

Build for current platform

```
gom build
```

## Release

```
$ rm -rf pkg/
$ mkdir -p pkg/
$ vendor/bin/gox -output="pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"
$ export VERSION=[v0.0.1]
$ vendor/bin/ghr $VERSION pkg/
$ git tag $VERSION
$ git push origin --tags
```
