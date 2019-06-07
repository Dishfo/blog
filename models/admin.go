package models

type Administrator struct {
	Id       int64  `gorm:"primary_key auto_increment"`
	Name     string `gorm:"unique_index;type:varchar(35)"`
	Password string `gorm:"type:varchar(35)"`
}

func (Administrator) TableName() string {
	return adminTable
}

const (
	adminTable = "administrator"
)

func QueryAdministrator(userName string) *Administrator {
	a := new(Administrator)
	dbInstance.Where("name = ?", userName).First(a)
	if a.Id == 0 {
		return nil
	}
	return a
}
