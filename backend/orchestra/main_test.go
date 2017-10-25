package orchestra

import (
	"reflect"
	"testing"
)

type testImageStream struct {
	stream   chan [][]uint8
	lastPath string
}

func TestProcessProducts(t *testing.T) {
	productsCh := make(chan ProductMessage, 1)
	imageStream := testImageStream{stream: make(chan [][]uint8, 1)}
	var imageCh chan ImageMessage
	productsCh <- ProductMessage{Path: "some path"}
	imageStream.stream <- [][]uint8{{1, 1, 1}}
	close(imageStream.stream)
	Process(productsCh, imageCh, &imageStream)
	if len(products) != 1 {
		t.Fatal("expected product to be appended to products")
	}
	if products[0].r != 1 ||
		products[0].g != 1 ||
		products[0].b != 1 {
		t.Fatal("unexpected values in product")
	}
	if imageStream.lastPath != "some path" {
		t.Fatal("fetched incorrect path", imageStream.lastPath)
	}
}

func TestFilterStream(t *testing.T) {
	ch := make(chan [][]uint8, 2)
	ch <- [][]uint8{{1, 1, 1}, {2, 0, 2}, {0, 0, 0}, {1, 1, 0}}
	ch <- [][]uint8{{0, 0, 0}, {2, 2, 2}, {0, 0, 0}, {1, 1, 1}}
	close(ch)
	retCh := FilterStream(ch)
	arr1 := <-retCh
	arr2 := <-retCh
	if !reflect.DeepEqual(arr1, [][]uint8{{1, 1, 1}, {2, 0, 2}, {1, 1, 0}}) {
		t.Fatal("incorrect values received from channel", arr1)
	}
	if !reflect.DeepEqual(arr2, [][]uint8{{2, 2, 2}, {1, 1, 1}}) {
		t.Fatal("incorrect values received from channel", arr2)
	}
}

func TestProcessImages(t *testing.T) {
	products = []Product{{1, 1, 1}, {0, 0, 0}}
	var productsCh chan ProductMessage
	stream := make(chan [][]uint8, 1)
	imageCh := make(chan ImageMessage, 1)
	imageCh <- ImageMessage{Path: "some path"}
	stream <- [][]uint8{{1, 1, 1}}
	imageStream := testImageStream{stream: stream}
	close(stream)
	Process(productsCh, imageCh, &imageStream)
	if len(matches) != 1 {
		t.Fatal("expected product to be appended to matches")
	}
	if len(matches[0].products) != len(products) {
		t.Fatal("unexpected products found in match")
	}

	//sort ok?
	m := matches[0]
	if m.products[1].r != 0 ||
		m.products[1].g != 0 ||
		m.products[1].b != 0 {
		t.Fatal("unexpected match found", m)
	}
	if m.products[0].r != 1 ||
		m.products[0].g != 1 ||
		m.products[0].b != 1 {
		t.Fatal("unexpected match found", m)
	}

	if imageStream.lastPath != "some path" {
		t.Fatal("fetched incorrect path", imageStream.lastPath)
	}
}

func (s *testImageStream) GetStream(path string) chan [][]uint8 {
	s.lastPath = path
	return s.stream
}
