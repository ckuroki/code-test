FROM golang:1.12.1-alpine
RUN apk add --no-cache make git curl
RUN mkdir -p /opt/shorten_url/

COPY config /opt/shorten_url/config
COPY config/docker-start.sh /opt/shorten_url

# Outside GOPATH to avoid issues using go modules
WORKDIR /src
COPY . .
RUN make 
RUN make install
EXPOSE 8080/tcp
ENTRYPOINT ["/opt/storage/docker-start.sh"]
