FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/bubble/frontend
RUN go get github.com/streadway/amqp
RUN go get github.com/astaxie/beego
RUN go get github.com/abbot/go-http-auth
CMD go install github.com/wwgberlin/bubble/frontend && /go/bin/frontend
EXPOSE 8080