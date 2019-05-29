package controllers

import (
	"blogServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
	"strings"
)

/**
对访问者的读写权限进行控制
*/
const (
	tokenPrefix = "Bearer"
)

var (
	limitPath = []string{
		"/v1/article/delete/:id",
		"/v1/article/edit",
		"/v1/article/add",
		"/v1/tag/add",
		"/v1/tag/delete/:id",
	}
)

func FilterUserPermission() {
	for _, url := range limitPath {
		beego.InsertFilter(url, beego.BeforeExec, checkUserPermission)
	}
}

func checkUserPermission(ctx *context.Context) {
	beego.BeeLogger.Info("check user permission")
	auths := ctx.Request.Header["Authorization"]
	for _, auth := range auths {
		if strings.HasPrefix(auth, tokenPrefix) {
			token := strings.Split(auth, " ")
			if len(token) > 0 && models.CheckToken(token[1]) {
				return

			}
		}
	}

	ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
	ctx.ResponseWriter.Flush()
}

//FilterIpAccessLimit 限制某一ip的访问频率
func FilterIpAccessLimit() {
	beego.InsertFilter("*", beego.BeforeExec, checkUserAccessLimit)
}

func checkUserAccessLimit(ctx *context.Context) {
	caddr := ctx.Request.RemoteAddr //ip:port
	caddrSpec := strings.Split(caddr, ":")
	host := caddrSpec[0]

	canLook := models.UserAccessSite(host)
	beego.BeeLogger.Info("%v", canLook)
	if !canLook {
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		ctx.ResponseWriter.Flush()
	}
}
