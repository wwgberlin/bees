FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/bubble/files-server
RUN go get github.com/streadway/amqp
CMD go install github.com/wwgberlin/bubble/files-server && /go/bin/files-server
EXPOSE 8080