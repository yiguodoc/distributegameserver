package main

import (
	_ "distributionGame/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.SetStaticPath("/javascripts", "static/js")
	beego.SetStaticPath("/bootstrap", "static/bootstrap")
	beego.SetStaticPath("/images", "static/img")
	beego.SetStaticPath("/dataTable", "static/dataTable")
	beego.SetStaticPath("/stylesheets", "static/css")
	beego.SetStaticPath("/framework7", "static/Framework7")

	beego.Run()
}
