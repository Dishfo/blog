package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//RegisterDB will be called in main package
func RegisterDB() {
	user := beego.AppConfig.String("mysql.user")
	pass := beego.AppConfig.String("mysql.pass")
	db := beego.AppConfig.String("mysql.db")
	host := beego.AppConfig.String("mysql.host")
	port := beego.AppConfig.String("mysql.port")

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pass,
		host, port, db)

	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", url)
	orm.RegisterModel(new(Tag), new(Article), new(Administrator))
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		beego.BeeLogger.Error("%s", err.Error())
	}
}
