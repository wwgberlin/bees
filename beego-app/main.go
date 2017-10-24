package main

import (
	"github.com/astaxie/beego"
	"github.com/wwgberlin/bubble/beego-app/controllers"
	"os"
)

func main() {
	os.MkdirAll("./images/", os.ModePerm)
	beego.Router("/", &controllers.UsersController{})
	beego.Run()
}
