// @APIVersion 1.0.0
// @Title web API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact 1771334691@qq.com
package routers

import (
	"blogServer/controllers"
	"github.com/astaxie/beego"
)

func init() {

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/admin",
			beego.NSInclude(
				&controllers.AdminController{},
			),
		),
		beego.NSNamespace("/tag",
			beego.NSInclude(
				&controllers.TagController{},
			),
		),
		beego.NSNamespace("/article",
			beego.NSInclude(&controllers.ArticleController{}),
		),
	)

	beego.AddNamespace(ns)
}
