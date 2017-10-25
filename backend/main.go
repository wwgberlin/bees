package main

import (
	"encoding/json"
	"github.com/wwgberlin/bubble/amqp"
	"github.com/wwgberlin/bubble/backend/orchestra"
)

type (
	productsFunc func(chan orchestra.ProductMessage)
	imagesFunc   func(chan orchestra.ImageMessage)
	conConsumer  struct {
		consumeNewImagesFunc   imagesFunc
		consumeNewProductsFunc productsFunc
	}
)

func main() {
	products := make(chan orchestra.ProductMessage, 10)
	products <- orchestra.ProductMessage{
		Path: "https://s-media-cache-ak0.pinimg.com/236x/4b/4c/55/4b4c559a7f7aa5ffdaec56c66efac0d6--beauty-pie-face-foundation.jpg",
	}
	products <- orchestra.ProductMessage{
		Path: "https://s-media-cache-ak0.pinimg.com/236x/3f/d0/f8/3fd0f84ee5df650222cc145e4d4c2ed3--beauty-pie-foundation.jpg",
	}
	products <- orchestra.ProductMessage{
		Path: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRAZmnui2ROTriNXeUGfQD7uRJMJjpH5RAgeLDiPtYYqABwP8nP7g",
	}
	products <- orchestra.ProductMessage{
		Path: "https://s-media-cache-ak0.pinimg.com/originals/6e/94/63/6e9463a2c3510854d6f11bb37e5ac03f.png",
	}

	c := conConsumer{
		consumeNewImagesFunc:   imagesFunc(consumeNewImages),
		consumeNewProductsFunc: productsFunc(consumeNewProducts(products)),
	}
	orchestra.Run(c, orchestra.NewImageStream())
}

func (c conConsumer) ConsumeNewImages(imagesCh chan orchestra.ImageMessage) {
	c.consumeNewImagesFunc(imagesCh)
}
func (c conConsumer) ConsumeNewProducts(productsCh chan orchestra.ProductMessage) {
	c.consumeNewProductsFunc(productsCh)
}

func consumeNewImages(imagesCh chan orchestra.ImageMessage) {
	messagesCh := make(chan amqp.Message)
	go consumeNewImagesMessage(messagesCh)
	for message := range messagesCh {
		var newImage orchestra.ImageMessage
		json.Unmarshal(message.Body(), &newImage)
		imagesCh <- newImage
		message.Ack()
	}
}
func consumeNewProducts(ch chan orchestra.ProductMessage) productsFunc {
	return func(productsCh chan orchestra.ProductMessage) {
		for {
			productsCh <- <-ch
		}
	}
}

func consumeNewImagesMessage(imagesCh chan amqp.Message) {
	config := amqp.GetChannelConfig()
	config.Name = "images"
	amqp.NewConsumerAMQPChannel(amqp.ConsumerChannelConfig{
		ChannelConfig: config,
		PrefetchCount: 1,
	}).Listen(imagesCh)
}
