package main

import (
	"encoding/json"
	"fmt"
	"github.com/wwgberlin/bubble/amqp"
)

func main() {
	conf := amqp.ConsumerChannelConfig{
		ChannelConfig: amqp.ChannelConfig{
			Name:     "jobs",
			User:     "guest",
			Password: "guest",
			Host:     "localhost",
			Port:     "5672",
		},
		PrefetchCount: 1,
	}
	amqp.NewConsumerAMQPChannel(conf)

}
