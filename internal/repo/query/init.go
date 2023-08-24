package query

import (
	"fmt"
	"github.com/Aliothmoon/Continu/internal/repo/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const MysqlDsn = "root:admin@(localhost:3306)/ci?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
	mysqlInit()
}
func mysqlInit() {
	db, err := gorm.Open(mysql.Open(MysqlDsn))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	SetDefault(db)
}

func sqliteInit() {
	db, err := gorm.Open(sqlite.Open("test.db"))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	err = db.AutoMigrate(&model.Log{}, &model.BuildRecord{}, &model.Project{})
	if err != nil {
		panic(err)
	}
	SetDefault(db)
}
