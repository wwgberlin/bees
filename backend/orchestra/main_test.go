package orchestra

import (
	"testing"
)

type (
	imagesFunc   func(chan ImageMessage)
	productsFunc func(chan ProductMessage)
	fakeConsumer struct {
		images                 chan ImageMessage
		products               chan ProductMessage
		consumeNewImagesFunc   imagesFunc
		consumeNewProductsFunc productsFunc
	}
	fakeImageStream struct {
		imageStream
	}
)

func TestRun(t *testing.T) {
	f:=fakeConsumer{}
	go Run(f)



	f.consumeNewImagesFunc = func(ch chan ImageMessage) {
		ch <- ImageMessage{}
	}
	f.consumeNewProductsFunc = func(ch chan ProductMessage) {
		ch <- ProductMessage{}
	}

}

func (f fakeConsumer) ConsumeNewImages(ch chan ImageMessage) {
	f.images = ch
}
func (f fakeConsumer) ConsumeNewProducts(ch chan ProductMessage) {
	f.products = ch
}
