FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/bubble/backend
RUN go get github.com/streadway/amqp
CMD go install github.com/wwgberlin/bubble/backend && /go/bin/backend
EXPOSE 8080