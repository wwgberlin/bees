package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abbot/go-http-auth"
	"github.com/astaxie/beego"
	"github.com/wwgberlin/bubble/amqp"
	"io"
)

type (
	User struct {
		Id           string
		PasswordHash string
	}
	DB []User
)

var db DB = DB([]User{
	{
		Id:           "wwgberlin",
		PasswordHash: "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1", // password is "hello"
	},
})

func Secret(user, realm string) string {
	for _, u := range db {
		if u.Id == user {
			return u.PasswordHash
		}
	}
	return ""
}

type UsersController struct {
	beego.Controller
}

func (c *UsersController) Prepare() {
	a := auth.NewBasicAuthenticator("example.com", Secret)
	if username := a.CheckAuth(c.Ctx.Request); username == "" {
		a.RequireAuth(c.Ctx.ResponseWriter, c.Ctx.Request)
	} else {
		c.Data["username"] = username
	}
}

func (c *UsersController) Get() {
	c.TplName = "upload.html"
}

func (c *UsersController) Post() {
	c.TplName = "post.html"
	file, header, err := c.GetFile("uploadfile")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// get the filename
	filename := header.Filename
	// save to server
	f, err := os.OpenFile("./images/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	b, _ := json.Marshal(struct {
		User interface{}
		Path string
	}{
		User: c.Data["username"],
		Path: "http://beego_app:8080/images/" + filename,
	})
	getPublisherChannel().Publish(b)
}

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
