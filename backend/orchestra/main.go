package orchestra

import (
	"fmt"
	"sort"
)

type Product struct {
	r uint8
	g uint8
	b uint8
}

type Consumer interface {
	ConsumeNewProducts(productsCh chan ProductMessage)
	ConsumeNewImages(productsCh chan ImageMessage)
}

var products []Product = []Product{}
var matches []ProductMatcher = []ProductMatcher{}

func Run(c Consumer, stream ImageStream) {
	imageCh := make(chan ImageMessage)
	productsCh := make(chan ProductMessage)
	go c.ConsumeNewImages(imageCh)
	go c.ConsumeNewProducts(productsCh)
	for {
		process(productsCh, imageCh, stream)
	}
}

func process(productsCh chan ProductMessage, imageCh chan ImageMessage, stream ImageStream) {
	select {
	case newProduct := <-productsCh:
		r, g, b := averageColor(filterStream(stream.GetStream(newProduct.Path)))
		products = append(products, Product{r: r, g: g, b: b})
		fmt.Println(products)
	case newImage := <-imageCh:
		r, g, b := averageColor(filterStream(stream.GetStream(newImage.Path)))
		m := newProductMatcher(products, r, g, b)
		sort.Sort(m)
		matches = append(matches, m)
		fmt.Println(matches)
	}
}

func averageColor(ch chan [][]uint8) (uint8, uint8, uint8) {
	var b, g, r uint64
	var size float64
	for arr := range ch {
		for i := range arr {
			b += uint64(arr[i][0])
			g += uint64(arr[i][1])
			r += uint64(arr[i][2])
			size++
		}
	}
	return uint8(float64(r) / size), uint8(float64(g) / size), uint8(float64(b) / size)
}

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
