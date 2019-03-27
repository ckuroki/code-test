# Shorten Url API

Compilation

* Install [Go](https://golang.org/doc/install#install)

* `make` to compile (will download dependencies if required, using [Go modules](https://github.com/golang/go/wiki/Modules) )

* `make test` to execute unit tests

* `make build-image` to build a docker container ( requires [Docker](http://docker.io) )

## Configuration and startup

1. Service configuration is stored into `config/config.yaml`

2. Start the API service with:

`./shorten_url_api`

