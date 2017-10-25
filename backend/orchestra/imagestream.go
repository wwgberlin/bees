package orchestra

import (
	"encoding/json"
	"io"
	"net/http"
)

func NewImageStream() ImageStream {
	return imageStream{}
}

type ImageStream interface {
	GetStream(path string) chan [][]uint8
}
type imageStream struct {
}

func (s imageStream) GetStream(path string) chan [][]uint8 {
	ch := make(chan [][]uint8)
	go func() {
		fetchImageSkinColor(path, ch)
		close(ch)
	}()
	return ch
}

func fetchImageSkinColor(path string, ch chan [][]uint8) {
	body := MustGet("http://cv_api:8080/api/skin/?image=" + path)
	defer body.Close()
	dec := json.NewDecoder(body)

	if _, err := dec.Token(); err != nil { //get rid of opening brackets
		panic(err)
	} else {
		for dec.More() {
			var arr [][]uint8
			if err := dec.Decode(&arr); err != nil {
				panic(err)
			}
			ch <- arr
		}
		if _, err = dec.Token(); err != nil { //get rid of closing brackets
			panic(err)
		}
	}
}

func MustGet(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp.Body
}
