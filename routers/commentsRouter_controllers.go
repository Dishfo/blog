package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["blogServer/controllers:AdminController"] = append(beego.GlobalControllerRouter["blogServer/controllers:AdminController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "AddArticle",
			Router:           `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "RemoveArticle",
			Router:           `/delete/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "EditArticle",
			Router:           `/edit`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "GetArticle",
			Router:           `/getArticle/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "GetArticleByTag",
			Router:           `/getArticleByTag`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "GetTopArticles",
			Router:           `/getTopArticles`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:ArticleController"] = append(beego.GlobalControllerRouter["blogServer/controllers:ArticleController"],
		beego.ControllerComments{
			Method:           "GetArticleList",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:TagController"] = append(beego.GlobalControllerRouter["blogServer/controllers:TagController"],
		beego.ControllerComments{
			Method:           "AddTag",
			Router:           `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:TagController"] = append(beego.GlobalControllerRouter["blogServer/controllers:TagController"],
		beego.ControllerComments{
			Method:           "RemoveTag",
			Router:           `/delete/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["blogServer/controllers:TagController"] = append(beego.GlobalControllerRouter["blogServer/controllers:TagController"],
		beego.ControllerComments{
			Method:           "GetAllTags",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
