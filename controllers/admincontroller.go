package controllers

import (
	"blogServer/models"
	"github.com/astaxie/beego"
	"net/http"
)

type AdminController struct {
	beego.Controller
}

func (c *AdminController) URLMapping() {
	c.Mapping("Login", c.Login)
}

// @Param user formData string true  "username"
// @Param pass formData string true  "password"
// @Success 200 {object} models.LoginResult
// @router /login [post]
func (c *AdminController) Login() {
	user := c.Ctx.Request.Form.Get("user")
	pass := c.Ctx.Request.Form.Get("pass")
	a := models.QueryAdministrator(user)
	if a != nil {
		if pass == a.Password {
			lr := models.LoginResult{}
			lr.Token = models.GenerateToken(user)
			lr.OperationResult = models.NewOperationResult(models.SUCCEED)
			c.Data["json"] = lr
			c.ServeJSON()
			return
		}
	}

	c.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
}
