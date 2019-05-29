package controllers

import "github.com/astaxie/beego"

type ErrController struct {
	beego.Controller
}

func (c *ErrController) Error401() {
	c.Ctx.ResponseWriter.WriteHeader(401)
	c.Ctx.WriteString("you need login for getting a token ")
}
