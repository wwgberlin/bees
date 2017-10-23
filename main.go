package main

import (
	"fmt"
	"github.com/parkr/imgix-go"
	"net/url"
)

func main() {
	client := imgix.NewClient("mycompany.imgix.net")

	// Simplest example
	fmt.Println(client.Path("/myImage.jpg"))

	// Example with parameters
	fmt.Println(client.PathWithParams("/myImage.jpg", url.Values{
		"w": []string{"400"},
		"h": []string{"300"},
	}))
}
