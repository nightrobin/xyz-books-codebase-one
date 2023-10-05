package method

import (
	"xyz-books/orm"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	Db = orm.ConnectToDB()
}
