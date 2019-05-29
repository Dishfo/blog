package models

import "github.com/astaxie/beego/orm"

type Administrator struct {
	Id       int64
	Name     string `orm:"unique"`
	Password string
}

const (
	adminTable = "administrator"
)

func QueryAdministrator(userName string) *Administrator {
	a := new(Administrator)
	o := orm.NewOrm()
	err := o.QueryTable(adminTable).Filter("Name", userName).One(a)
	if err != nil {
		return nil
	}
	if a.Id == 0 {
		return nil
	}
	return a
}
