package orchestra

import (
	"encoding/json"
	"io"
	"net/http"
)

type (
	ImageMessage struct {
		User interface{}
		Path string
	}
)

func filterStream(ch chan [][]uint8) chan [][]uint8 {
	newCh := make(chan [][]uint8)
	go func() {
		for arr := range ch {
			newArr := [][]uint8{}
			for i := range arr {
				if arr[i][0] != 0 ||
					arr[i][1] != 0 ||
					arr[i][2] != 0 {
					newArr = append(newArr, arr[i])
				}
			}
			newCh <- newArr
		}
		close(newCh)
	}()
	return newCh
}

func imageStream(newImage ImageMessage) chan [][]uint8 {
	ch := make(chan [][]uint8)
	go func() {
		fetchImageFace(newImage, ch)
		close(ch)
	}()
	return ch
}

func fetchImageFace(newImage ImageMessage, ch chan [][]uint8) {
	body := MustGet("http://cv_api:8080/api/skin/?image=" + newImage.Path)
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
