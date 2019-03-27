PRJ_DIR		?= $(shell pwd)

all: shorten_url_api

shorten_url_api: main.go
	go build -o shorten_url_api main.go

build-image: 
	docker build -t  "shorten_url_api:latest"  .
	docker run --rm -w /source -v ${PRJ_DIR}:/source "shorten_url_api:latest" 

test:
	go test -v . 

lint:
	golint . 

clean: 
	rm -f shorten_url_api 

install: shorten_url_api 
	cp shorten_url_api /opt/shorten_url
	cp config/docker-start.sh /opt/shorten_url
	chmod +x /opt/storage/docker-start.sh

setup:
	go get github.com/aws/aws-sdk-go
	go get github.com/gin-gonic/gin
	go get github.com/gin-gonic/gin
	go get github.com/itsjamie/gin-cors
	go get github.com/sirupsen/logrus
	go get github.com/spf13/viper
	go get github.com/appleboy/gofight
	go get github.com/stretchr/testify/assert
