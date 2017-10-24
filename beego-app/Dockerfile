FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/bubble/beego-app
RUN go get github.com/streadway/amqp
RUN go get github.com/astaxie/beego
RUN go get github.com/abbot/go-http-auth
CMD go install github.com/wwgberlin/bubble/beego-app && /go/bin/beego-app
EXPOSE 8080