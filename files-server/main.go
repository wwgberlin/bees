package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"encoding/json"
	"github.com/wwgberlin/bubble/amqp"
)

func getPublisherChannel() amqp.PublisherChannel {
	file, err := os.Open("../config/rabbitmq.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	configuration := amqp.ChannelConfig{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return amqp.NewPublisherChannel(amqp.PublisherChannelConfig{
		ChannelConfig: configuration,
		Exchange:      "images",
	})
}

func main() {
	os.MkdirAll("./static/images/", os.ModePerm)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", upload(os.Getenv("HOSTNAME"), getPublisherChannel()))
	http.ListenAndServe(":8080", nil)
}

func upload(domain string, ch amqp.PublisherChannel) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		filename := handler.Filename
		f, err := os.OpenFile("./static/images/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		b, _ := json.Marshal(struct {
			Path string
		}{"http://" + domain + ":8080/static/images/" + filename})
		ch.Publish(b)
		fmt.Fprint(w, "all good!")
	}
}
