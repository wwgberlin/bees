package orchestra

import (
	"fmt"
	"sort"
	"time"
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
	go func() {
		for {
			fmt.Println("products:", products)
			fmt.Println("matches:", matches)
			fmt.Println("=================")
			time.Sleep(time.Second * 10)
		}
	}()
	for {
		Process(productsCh, imageCh, stream)
	}
}

/**
use select in this function to sync processing of new images and processing of new products
your helper methods:
	* use stream.GetStream(image path) to get an image with skin color
	* use FilterStream(chan [][]uint8) to get rid of all the black RGB values given by the skin detector
	* use averageColor(chan [][]uint8) to compute the average color of the image

On new product:
	* get the average RBG value of the product
	* append a new product instance with the average values

On new user image:
	* get the average RBG value of the image
	* create a new productMatcher with the products and the RBG colors.
	* sort the productMatcher
	* append the productMatcher to the average values
*/

func Process(productsCh chan ProductMessage, imageCh chan ImageMessage, stream ImageStream) {
	select {
	case newProduct := <-productsCh:
		r, g, b := averageColor(FilterStream(stream.GetStream(newProduct.Path)))
		products = append(products, Product{r: r, g: g, b: b})
	case newImage := <-imageCh:
		r, g, b := averageColor(FilterStream(stream.GetStream(newImage.Path)))
		m := newProductMatcher(products, r, g, b)
		sort.Sort(m)
		matches = append(matches, m)
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

/** FilterStream iterates over the values in a channel and shoves them to a new channel
items will be pushed to the channel only if their values are not RGB values 0,0,0
* append is your friend.
* remember to close the channel when the channel the function is listening to closes
*/
func FilterStream(ch chan [][]uint8) chan [][]uint8 {
	newCh := make(chan [][]uint8)
	go func() {
		for arr := range ch {
			newArr := [][]uint8{}
			for i := range arr {
				if arr[i][0] != 0 ||
					arr[i][1] != 0 ||
					arr[i][2] != 0 {
					newArr = append(newArr, []uint8{arr[i][0], arr[i][1], arr[i][2]})
				}
			}
			fmt.Println("<<<<<<<", newArr)
			newCh <- newArr
		}
		fmt.Println("closing")
		close(newCh)
	}()
	return newCh
}
