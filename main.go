package main

import (
	"blogServer/controllers"
	"blogServer/models"
	_ "blogServer/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func main() {
	beego.ErrorController(&controllers.ErrController{})
	models.InitModels()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
			ctx.Output.Header("Access-Control-Allow-Origin", "*")
		})
	}

	if beego.BConfig.RunMode == "prod" {
		controllers.FilterUserPermission()
		controllers.FilterIpAccessLimit()
	}
	/*	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))*/
	beego.Run()
}
