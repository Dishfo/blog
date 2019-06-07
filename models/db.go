package models

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	dbInstance *gorm.DB
)

//RegisterDB will be called in main package
func RegisterDB() {
	var err error
	user := beego.AppConfig.String("mysql.user")
	pass := beego.AppConfig.String("mysql.pass")
	dbName := beego.AppConfig.String("mysql.db")
	host := beego.AppConfig.String("mysql.host")
	port := beego.AppConfig.String("mysql.port")

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, pass,
		host, port, dbName)

	dbInstance, err = gorm.Open("mysql", url)
	if err != nil {
		log.Fatal(err.Error())
	}

	errors :=
		dbInstance.AutoMigrate(new(Tag), new(Article), new(Administrator)).GetErrors()

	for _, err := range errors {
		beego.BeeLogger.Error("%s", err.Error())
	}

	dbInstance.DB().SetMaxIdleConns(20)
	dbInstance.LogMode(true)
}
