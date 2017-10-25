package orchestra_go

import (
	"encoding/json"
	"fmt"
	"github.com/wwgberlin/bubble/amqp"
	"os"
)

func Run() {
	imageCh := make(chan amqp.Message)
	go consumeNewImages(imageCh)
	select {
	case message := <-imageCh:
		var newImage ImageMessage
		json.Unmarshal(message.Body(), &newImage)
		var b, g, r uint64
		var size float64
		for arr := range filterStream(imageStream(newImage)) {
			for i := range arr {
				b += uint64(arr[i][0])
				g += uint64(arr[i][1])
				r += uint64(arr[i][2])
				size++
			}
		}
		message.Ack()
		fmt.Println("User", newImage.User, "skin color is [", float64(r)/size, float64(g)/size, float64(b)/size, "]")
	}
}

func consumeNewImages(imagesCh chan amqp.Message) {
	file, err := os.Open("../config/rabbitmq.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	config := amqp.ChannelConfig{}
	err = decoder.Decode(&config)
	config.Name = "images"
	if err != nil {
		fmt.Println("error:", err)
	}
	amqp.NewConsumerAMQPChannel(amqp.ConsumerChannelConfig{
		ChannelConfig: config,
		PrefetchCount: 1,
	}).Listen(imagesCh)
}
